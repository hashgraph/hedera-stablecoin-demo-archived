#!/bin/sh

git pull
./mvnw install
cp target/*.jar runme.jar
