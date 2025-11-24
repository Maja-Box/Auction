# Auction

Program by Madeleine, Mads and Vee

Welcome to our program, Auction.

Description:

This program simulates an Auction system where clients can bid and the servers can handle crashes up to 2 times. 

How to use:

Run the server program three times with numbers 0 to 2 (will be prompted upon server start) to choose a port
for each server.

Thereafter, start as many clients as you want (each in a seperate terminal)

Then type in a number for the bid.

thereafter you can type bid or check to either bid again or check the current state of the auction. 

To crash a server, go into the server 5050 terminal and press CTRL+C. 

After a crash, when client tries to bid again, it will result in a error where client will change server.

Known issues:

For some reason, after crashing the first server. Whenever a client tries to bid again or check the status
of the auction, a nil value will appear after the normal message. We have no clue how to fix this issue but it isnt causing any fundamental errors. 

# Logs

<img width="851" height="154" alt="image" src="https://github.com/user-attachments/assets/7e55d693-929f-453e-9c45-bd871a33408a" />

<img width="854" height="161" alt="image" src="https://github.com/user-attachments/assets/e9db52f2-5e60-4923-abc3-1d768a8e1ebe" />

<img width="868" height="147" alt="image" src="https://github.com/user-attachments/assets/5a4cf354-b8c8-4fca-a4d5-15b01735716d" />

<img width="1364" height="737" alt="image" src="https://github.com/user-attachments/assets/9eaaefe1-e887-4a00-bad1-564e389f547e" />

