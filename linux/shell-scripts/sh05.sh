#!/bin/bash
case $1 in
    "hello")
      echo "Hello world
      ;;
    "")
      echo "no world input"
    ;;
    *)
      echo "* was trigger"
    ;;
  esac
exit 0
