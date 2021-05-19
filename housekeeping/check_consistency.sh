#!/bin/bash

grep -r fmt.Printf "$1"
grep -r fmt.Println "$1"
grep -r writeString "$1" 
grep -r bufio.NewWriter "$1"
grep -r bufio.NewReader "$1"
