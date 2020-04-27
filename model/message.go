/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 26-04-2020
 * |
 * | File Name:     message.go
 * +===============================================
 */

package model

import "time"

// Message represents a message to broadcast
type Message struct {
	From      string
	CreatedAt time.Time
	Text      string
}
