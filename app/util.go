// Copyright 2020 The Lunar.Music.Web AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getTimeDisplay(millisecond int) string {
	hour := millisecond / 1000 / 60 / 60
	minute := millisecond/1000/60 - hour*60
	second := millisecond/1000 - hour*60*60 - minute*60

	if hour == 0 {
		return fmt.Sprintf("%02d", minute) + ":" + fmt.Sprintf("%02d", second)
	}

	return fmt.Sprintf("%02d", hour) + ":" + fmt.Sprintf("%02d", minute) + ":" + fmt.Sprintf("%02d", second)
}
