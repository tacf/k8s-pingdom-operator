/*
Copyright 2019 github.com/tacf.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pingdom

import (
	"fmt"

	pingdom "github.com/russellcardullo/go-pingdom/pingdom"
	"k8s.io/klog"
)

// Credentials holds needed values for single user auth
type Credentials struct {
	Username string
	Password string
	Apikey   string
}

// BasicHTTPCheck represents a data type holding the needed
// details for setting up a Http check
type BasicHTTPCheck struct {
	Name     string
	URL      string
	Interval uint8
}

// Client is a wrapper type around go-pingdom clinet
type Client pingdom.Client

// NewClient creates single user pingdom client
func NewClient(credentials Credentials) (Client, error) {
	client, err := pingdom.NewClientWithConfig(pingdom.ClientConfig{
		User:     credentials.Username,
		Password: credentials.Password,
		APIKey:   credentials.Apikey,
	})

	return Client(*client), err // Converting to alias type
}

func (c *Client) getChecks() []pingdom.CheckResponse {
	pingdomChecks, _ := c.Checks.List()
	return pingdomChecks
}

func (c *Client) getCheckDetails(check BasicHTTPCheck) (pingdom.CheckResponse, error) {
	for _, pdCheck := range c.getChecks() {
		if pdCheck.Name == check.Name {
			return pdCheck, nil
		}
	}
	return pingdom.CheckResponse{}, fmt.Errorf("Could not retrieve check details for %s", check.Name)
}

func (c *Client) checkExists(check BasicHTTPCheck) bool {
	for _, pdCheck := range c.getChecks() {
		if pdCheck.Name == check.Name {
			return true
		}
	}
	return false
}

func (c *Client) checkNeedsUpdate(check BasicHTTPCheck) bool {
	checkDetails, err := c.getCheckDetails(check)
	if err != nil {
		return true // if check does not exist assume "update needed"
	}
	return !(checkDetails.Hostname == check.URL && checkDetails.Resolution == int(check.Interval))
}

func (c *Client) updateCheck(check BasicHTTPCheck) {
	if c.checkNeedsUpdate(check) {
		klog.Infof("Check %s matches definition found in Pingdom...Updating it!", check.Name)
		checkDetails, _ := c.getCheckDetails(check) // checkNeedsUpdate already validates check existance
		newDetails := pingdom.HttpCheck{Name: check.Name, Hostname: check.URL, Resolution: int(check.Interval)}
		_, err := c.Checks.Update(checkDetails.ID, &newDetails)
		if err != nil {
			klog.Errorf("Error updating check from resource %s", check.Name)
		} else {
			klog.Infof("Successfully updated check details with new resource %s", check.Name)
		}
	}
	return
}

func (c *Client) deleteCheck(check BasicHTTPCheck) {
	if c.checkExists(check) {
		klog.Infof("Check %s matches definition found in Pingdom ... Deleting it!", check.Name)
		checkDetails, err := c.getCheckDetails(check)
		if err != nil {
			klog.Errorf("Check %s - %s not found", check.Name, check.URL)
			return
		}
		_, err = c.Checks.Delete(checkDetails.ID)
		if err != nil {
			klog.Errorf("Error deleting check from resource %s", check.Name)
		} else {
			klog.Infof("Successfully deleted check %s", check.Name)
		}
	} else {
		klog.Warningf("Check from resource %s not found in Pingdom ... Skipping it!")
	}
	return
}

// AddCheck add a new check to pingdom service
func (c *Client) AddCheck(checks []BasicHTTPCheck) {
	for _, check := range checks {
		klog.Infof("PINGDOM ADD -> Url: %s, Interval:%d", check.URL, check.Interval)

		// Try to sync with existing checks
		if c.checkExists(check) {
			klog.Infof("Check %s already exists", check.Name)
			c.updateCheck(check)
			return
		}

		newCheck := pingdom.HttpCheck{Name: check.Name, Hostname: check.URL, Resolution: int(check.Interval)}
		_, err := c.Checks.Create(&newCheck)
		if err != nil {
			klog.Errorf("Could not create check: %s", err)
			return
		}
		klog.Infof("Created check: %s, URL: %s, Interval: %d ", check.Name, check.URL, check.Interval)
	}
}

// UpdateCheck updates existing check on pingdom service
func (c *Client) UpdateCheck(checks []BasicHTTPCheck) {
	for _, check := range checks {
		klog.Infof("PINGDOM UPDATE -> Url: %s, Interval:%d", check.URL, check.Interval)
		c.updateCheck(check)
	}
}

// DeleteCheck deletes existing check on pingdom service
func (c *Client) DeleteCheck(checks []BasicHTTPCheck) {
	for _, check := range checks {
		klog.Infof("PINGDOM DELETE -> Url: %s, Interval:%d", check.URL, check.Interval)
		c.deleteCheck(check)
	}
}
