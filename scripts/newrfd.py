#import argparse
import os
import sys
import re

class bcolors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKCYAN = '\033[96m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'

print("======================================================================")
print("CREATE NEW RFD")
print("======================================================================")

stream = os.popen('git config user.name')
user_name = stream.read()
print('Hello {}'.format(user_name))


rfd_title = raw_input("Please enter title for the RFD:\n")
print("")

state_ok = False
state_str = "ideation"

while (not(state_ok)):

    print("Please select an initial state for the RFD")
    print("     1. Ideation")
    print("     2. Prediscussion")
    print("")
    val = raw_input(" Enter 1 or 2:")

    print("")

    try:
        state_ok = True
        state = int(val)
        if (state == 1):
            print("You selected ideation")
        
        elif (state == 2):
            state_str = "prediscussion"
            print("You selected prediscussion")
        
        else:
            raise ValueError
    
    except ValueError:
        print("Please select 1 or 2")
        state_ok = False

print("The RFD will be created with the following metadata:")

print(bcolors.OKGREEN)
print("---")
print("title: {}".format(rfd_title))
print("authors: {}".format(user_name))
print("state: {}".format(state_str))
print("discussion: <link to discussion")
print("---")
print(bcolors.ENDC)

cont = raw_input("Continue? Y/n\n")

if (cont == "n"):
    sys.exit()

if (not(cont == "Y" or cont == "y")):
    print("Invalid input ... exiting")
    sys.exit()

# First check repository for latest branch no


