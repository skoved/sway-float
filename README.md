# sway-float
Simple program to float windows immediately.

## Problem
Firefox based browsers (and other programs) create a new window and immediately give the window a default name
(ex: Mozilla Firefox) before setting the actual name of the window. Sway has trouble dealing with this change in window
title names and will only match the first title that is set until the window focus changes. Then the it will match the
new title and the window will be made floating. This program uses sway-ipc to listen directly for window events and will
set a window to floating when a new window that matches the criteria is created.
