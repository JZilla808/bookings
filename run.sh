#!/bin/bash

go build -o bookings cmd/web/*.go && ./bookings -cache=false -production=false