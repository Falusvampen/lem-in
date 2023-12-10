# The purpose of this file is to convert a file from LF to CRLF
# Usage: python LFtoCRLF.py

import glob
import os

# get a list of all the files in the current directory
files = glob.glob("*")

# go through all the files and replace the /n with /r/n
for file in files:
    # skip files and directories that start with "." or are named "LFtoCRLF.py"
    if not file.startswith(".") and file != "LFtoCRLF.py":
        print("processing file: ", file)
        # open the file for reading
        f = open(file, "r")
        # read the contents of the file
        contents = f.read()
        # close the file
        f.close()
        
        # reopen the file for writing
        f = open(file, "w")
        # change the line endings
        contents = contents.replace("\n", "\r\n")
        # write the contents back to the file
        f.seek(0)
        f.write(contents)
        f.close()
