#!/bin/sh
# if go.mod is not in the current directory, create that file and it's content
FILE=./go.mod 
if test ! -f "$FILE"; then
  go mod init falusvampen-lem-in-try
  go mod tidy
fi

# Create colors to use
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[34m'
MAGENTA='\033[35m'
ORANGE='\033[33m'
CYAN='\033[36m'
NC='\033[0m' # No Color
CHECK_MARK="${GREEN}\xE2\x9C\x94${NC}"
#background colors
BDG='\033[100m' #BACK_DARK_GRAY
BLB='\033[104m' #BACK_Light blue
BACK_RED='\033[101m'

# background line seperator
sep() {
 sleep 0.2;
 echo; printf "${BLB}.${NC}"; echo;
}

# create a visual boundary of 50 dots
boundary(){
 for i in `seq 1 51`
 do
  if [ $(expr $i % 3) == 0 ]; then printf "${CYAN}.${NC}"
  elif [ $(expr $i % 2) == 0 ]; then printf "${RED}.${NC}"
  else echo -ne "."; fi
  sleep 0.04
 done
 echo
 echo
}

#congratulations message
congrats=(
" ${CHECK_MARK} ${GREEN}Congratulations! ${GREEN}All tests passed...${NC} ${CHECK_MARK} "
" ${CHECK_MARK} ${YELLOW}Congratulations! ${GREEN}All tests passed...${NC} ${CHECK_MARK} "
)

#go through congrats array
j=0;
congratulations(){
 while [ $j -le 10 ]
 do 
  for i in "${congrats[@]}"
  do
   echo -ne "\r$i"
   sleep 0.4
   j=$(($j+1))
  done
 done
}

#test files
testfiles=(
"badexample00.txt"
"badexample01.txt"
"example00.txt"
"example01.txt"
"example02.txt"
"example03.txt"
"example04.txt"
"example05.txt"
# "example06.txt"
# "example07.txt"
)

echo; sep; echo; printf ${CYAN}"Running tests..."${NC}; boundary; echo;

#iterate through test files
for i in "${testfiles[@]}"
  do
   printf ${BLUE}"go run . ${i}"${NC}; echo; echo; sleep 0.8;
   go run . ${i}; echo;
   sleep 2; echo; boundary; echo;
  done

echo 

congratulations
echo
echo

printf "${CYAN} 
                    o      o 
                    o  o   
                        o     o

                    o    o o 
                ________.____________
                |  .                |
                |^^^.^^^^^.^^^^^.^^^|
                |   .  .  .         |
                 \   . . . .       /
C H E E R S !!!   \   . .        / 
                    \  ..       / 
                     \        /
                      \     /
                       \  /
                        ||
                        ||
                        ||
                        ||
                        ||
                        /\\
                       /;;\\
                      /;;;;\\
         ==============================
         ||   ♪┏(・o･)┛♪┗ ( ･o･) ┓♪   ||
${GREEN} 
   ____            _____   _____
 | __  \   /\     / ____| / ____|
 | |__) | /  \   | (___  | (___  
 | | ___// /\ \   \___  \ \___  \ 
 | |    / ____ \  ____) |  ____) |
 |_|   /_/    \_\|_____/  |_____/  
 "

sep
echo 
