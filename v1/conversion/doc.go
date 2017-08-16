// This entire package is a short-term fix for a bigger problem of sharing
// conversion logic between Loggregator and its client.
//
// *DO NOT* make changes to this code here.
//
// For the time being, this code is a duplicate of Loggregator's conversion
// package and it should be updated regularly until we resolve:
//
// https://github.com/cloudfoundry-incubator/go-loggregator/issues/19
//
// Current copy of this package is from Loggregator Git SHA:
// 737384a568fa12cf9115e48afc75d191fdcc70bb
package conversion
