package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// The Graph structure keeps track of all rooms the ant can take, the start and end rooms of the path and the number of ants
type Graph struct {
	Rooms         []*Room
	StartRoomName string
	EndRoomName   string
	Ants          int
}

// Room represents a room in an ant hill
type FRoom struct {
	Name        string   // name of the room
	X           int      // x-coordinate of the room
	Y           int      // y-coordinate of the room
	Connections []*FRoom // list of rooms connected to this room
	Visited     bool     // whether this room has been visited or not
	Distance    int      // distance of this room from the start room
}

// AntHill represents an ant hill with rooms, a start room, an end room, and the number of ants
type AntHill struct {
	FRooms    []*FRoom // list of all rooms in the ant hill
	StartRoom *FRoom   // start room of the ant hill
	EndRoom   *FRoom   // end room of the ant hill
	Ants      int      // number of ants in the ant hill
}

// The Room structure keeps track of the roomname, The rooms that the the current room is connected to and if the room has been visited before
type Room struct {
	Roomname    string
	Connections []string
	Visited     bool
}

var ah AntHill // global variable representing the ant hill

// AddRoom is a method that adds a new room, name, to a graph
func (g *Graph) AddRoom(name string) {
	g.Rooms = append(g.Rooms, &Room{Roomname: name, Connections: []string{}, Visited: false})
}

// AddLinks is a method that adds a link from one room to another
func (g *Graph) AddLinks(from, to string) {
	fromRoom := g.getRoom(from)
	toRoom := g.getRoom(to)
	if fromRoom == nil || toRoom == nil {
		log.Fatalf("Room doesn't exist (%v-%v)", from, to)
	}
	if contains(fromRoom.Connections, to) || contains(toRoom.Connections, from) {
		log.Fatalf("ERROR: invalid data format. Duplicate Link (%v --- %v)", from, to)
	}
	switch {
	case fromRoom.Roomname == g.EndRoomName:
		toRoom.Connections = append(toRoom.Connections, fromRoom.Roomname)
	case toRoom.Roomname == g.EndRoomName:
		fromRoom.Connections = append(fromRoom.Connections, toRoom.Roomname)
	case toRoom.Roomname == g.StartRoomName:
		toRoom.Connections = append(toRoom.Connections, fromRoom.Roomname)
	case fromRoom.Roomname == g.StartRoomName:
		fromRoom.Connections = append(fromRoom.Connections, toRoom.Roomname)
	default:
		fromRoom.Connections = append(fromRoom.Connections, toRoom.Roomname)
		toRoom.Connections = append(toRoom.Connections, fromRoom.Roomname)
	}
}

func (g *Graph) getRoom(name string) *Room {
	for _, room := range g.Rooms {
		if room.Roomname == name {
			return room
		}
	}
	return nil
}

// ReadFile reads the given file and returns its contents as a list of lines
func ReadFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var originalFileLines []string
	for scanner.Scan() {
		originalFileLines = append(originalFileLines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return originalFileLines, nil
}

func main() {
	log.SetFlags(0)

	// ///////////////////////////////////////////////////
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <filename>")
		return
	}
	originalFileLines, err := ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Remove the comments from the original file lines
	filteredLines := RemoveComments(originalFileLines)

	// check length of slice to be minimum 6: 1st line is number of ants, 2nd  and 3rd line is start room, 4th and 5th line is end room, 6th line is a link
	if len(filteredLines) < 6 {
		NoGo("")
	}

	// check if first line is a number
	if !IsNumber(filteredLines[0]) {
		NoGo("")
	}

	// convert first line to int and store in AntNum
	ah.Ants, _ = strconv.Atoi(filteredLines[0])
	filteredLines = filteredLines[1:]

	// check if number of ants is valid
	if ah.Ants <= 0 {
		NoGo("Number of ants is invalid")
	}

	No2Dashes(filteredLines)
	No3Spaces(filteredLines)
	NoDuplicateLines(filteredLines)
	NoHashInLastLine(filteredLines)

	// extract start room
	ExtractStartRoom(filteredLines)
	filteredLines = DeleteStartRoom(filteredLines)

	// extract end room
	ExtractEndRoom(filteredLines)
	filteredLines = DeleteEndRoom(filteredLines)

	// extract rooms
	ExtractRooms(filteredLines)
	OnlyConnections := DeleteAllRooms(filteredLines)

	// check if any room is there in the connections that is not in the rooms
	CheckRoomsInConnectionsPresent(OnlyConnections, GetAllRoomNames(&ah))

	// Add Connections to the rooms where a connection is in the format "room1-room2" and room1 and room2 are in the rooms
	AddConnections(OnlyConnections)

	checkUnconnectedRooms(&ah)

	/////////////////////////////////////////////////

	lines := validateFileGiveMeStrings()
	_ = lines
	//create a graph for the DFS search
	gdfs := &Graph{Rooms: []*Room{}}
	if err := PopulateGraph(lines, gdfs); err != nil {
		fmt.Print(err)
		return
	}
	gbfs := DeepCopyGraph(gdfs)
	// Print the contents of the slice with a new line after each element
	fmt.Println(strings.Join(originalFileLines, "\n") + "\n")

	allPathsDFS, allPathsBFS := []string{}, []string{}
	var path string
	DFS(gdfs.StartRoomName, gdfs.EndRoomName, gdfs, path, &allPathsDFS)
	BFS(gbfs.StartRoomName, gbfs.EndRoomName, gbfs, &allPathsBFS, ShortestPath)
	lenSorter(&allPathsBFS)
	lenSorter(&allPathsDFS)
	antNum := gdfs.Ants
	DFSSearch := AntSender(antNum, allPathsDFS)
	BFSSearch := AntSender(antNum, allPathsBFS)

	if len(DFSSearch) == 0 || len(BFSSearch) == 0 {
		if len(DFSSearch) == 0 {
			fmt.Println(fmt.Errorf("DFS Search Failed").Error())
		}
		if len(BFSSearch) == 0 {
			fmt.Println(fmt.Errorf("BFS Search Failed").Error())
		}
		return
	}

	for _, step := range shorterSearch(DFSSearch, BFSSearch) {
		fmt.Println(step)
	}
}

// BFS preforms a Breadth First Search of a graph from rooms start to end and puts all paths found in the []string paths
func BFS(start, end string, g *Graph, paths *[]string, f func(graph *Graph, start string, end string, path []string) []string) {
	begin := g.getRoom(start)

	for i := 0; i < len(begin.Connections); i++ {
		var shortPath []string
		ShortestPath(g, g.StartRoomName, g.EndRoomName, shortPath)
		var shortStorer string
		if len(pathArray) != 0 {
			shortStorer = pathArray[0]
		}

		for _, v := range pathArray {
			if len(v) < len(shortStorer) {
				shortStorer = v
			}
		}

		if len(pathArray) != 0 {
			shortStorer = shortStorer[1 : len(shortStorer)-1]
		}

		shortStorerSlc := strings.Split(shortStorer, " ")
		shortStorerSlc = shortStorerSlc[1:]

		for z := 0; z < len(shortStorerSlc)-1; z++ {
			g.getRoom(shortStorerSlc[z]).Visited = true
		}

		var pathStr string
		if len(shortStorerSlc) != 0 {
			for i := 0; i < len(shortStorerSlc); i++ {
				if i == len(shortStorerSlc)-1 {
					pathStr += shortStorerSlc[i]
				} else {
					pathStr = pathStr + shortStorerSlc[i] + "-"
				}
			}
		}

		if len(pathStr) != 0 {
			if len(pathStr) != 0 {
				containing := false
				for _, v := range *paths {
					if v == pathStr {
						containing = true
					}
				}
				if !containing {
					*paths = append(*paths, pathStr)
				}
			}
			pathArray = []string{}
		}
	}
}

// PopulateGraph sorts goes through the txt file with the ants and rooms and adds the rooms and links to the graph
func PopulateGraph(lines []string, g *Graph) error {
	var err error
	// Parse the number of ants from the first line
	g.Ants, err = strconv.Atoi(lines[0])
	if err != nil {
		return err
	}
	if g.Ants == 0 {
		return errors.New("ERROR: invalid data format. Number of ants must be greater than 0")
	}

	// Iterate over the remaining lines to populate the graph
	start := false
	end := false
	for _, line := range lines[1:] {
		space := strings.Split(line, " ")

		if len(space) > 1 {
			if !isValidRoomName(line) {
				log.Fatalf("ERROR: invalid data format. Room name or room coordinates invalid")
			}
			g.AddRoom(space[0])
		}

		if start {
			g.StartRoomName = g.Rooms[len(g.Rooms)-1].Roomname
			start = false
		} else if end {
			g.EndRoomName = g.Rooms[len(g.Rooms)-1].Roomname
			end = false
		}

		hyphen := strings.Split(line, "-")
		if len(hyphen) > 1 {
			if hyphen[0] == hyphen[1] {
				log.Fatalf("ERROR: invalid data format. You have a connection from the same room to same room.\n")
			}
			g.AddLinks(hyphen[0], hyphen[1])
		}
		switch line {
		case "##start":
			start = true
		case "##end":
			end = true
		}
	}

	return nil
}

// DFS preforms a depth first search of a graph and returns the possible paths
func DFS(current, end string, g *Graph, path string, pathList *[]string) {
	curr := g.getRoom(current)
	if current != end {
		curr.Visited = true
	}
	if curr.Roomname == g.EndRoomName {
		path += current
	} else if !(curr.Roomname == g.StartRoomName) {
		path += current + "-"
	}

	if current == end {
		*pathList = append(*pathList, path)
		path = ""
		for i := 0; i < len(g.getRoom(g.StartRoomName).Connections); i++ {
			if g.getRoom(g.StartRoomName).Connections[i] == g.EndRoomName {
				g.getRoom(g.StartRoomName).Connections[i] = ""
			}
		}
		DFS(g.StartRoomName, end, g, path, pathList)
	}
	for i := 0; i < len(curr.Connections); i++ {
		if curr.Connections[i] == g.EndRoomName {
			curr.Connections[0], curr.Connections[i] = curr.Connections[i], curr.Connections[0]
		}
	}
	for _, roomName := range curr.Connections {
		if roomName == "" {
			continue
		}
		currRoom := g.getRoom(roomName)
		if !currRoom.Visited {
			DFS(currRoom.Roomname, end, g, path, pathList)
		}
	}
}

var pathArray []string

// ShortestPath finds all the possible paths from start to end room using BFS and sorts them in ascending order
func ShortestPath(graph *Graph, start string, end string, path []string) []string {
	path = append(path, start)
	if start == end {
		return path
	}
	shortest := make([]string, 0)
	for _, node := range graph.getRoom(start).Connections {
		if !contains(path, node) && !graph.isVisited(node) {
			newPath := ShortestPath(graph, node, end, path)
			if len(newPath) > 0 && contains(newPath, graph.StartRoomName) && contains(newPath, end) {
				pathArray = append(pathArray, fmt.Sprint(newPath))
			}
		}
	}
	return shortest
}

func (graph *Graph) isVisited(str string) bool {
	return graph.getRoom(str).Visited
}

func lenSorter(paths *[]string) {
	sort.Slice(*paths, func(i, j int) bool {
		return len((*paths)[i]) < len((*paths)[j])
	})
}

func AntSender(n int, pathList []string) []string {
	pathLists := make([][]string, len(pathList))
	for i, path := range pathList {
		pathLists[i] = strings.Split(path, "-")
	}

	queue := make([][]string, len(pathList))

	for i := 1; i <= n; i++ {
		minStepsIndex := 0
		minSteps := len(pathLists[0]) + len(queue[0])
		for j, path := range pathLists {
			steps := len(path) + len(queue[j])
			if steps < minSteps {
				minSteps = steps
				minStepsIndex = j
			}
		}
		queue[minStepsIndex] = append(queue[minStepsIndex], strconv.Itoa(i))
	}

	container := make([][][]string, len(queue))
	for i, path := range queue {
		for _, ant := range path {
			adder := make([]string, len(pathLists[i]))
			for j, room := range pathLists[i] {
				adder[j] = "L" + ant + "-" + room
			}
			container[i] = append(container[i], adder)
		}
	}

	finalMoves := []string{}
	for _, paths := range container {
		for j, moves := range paths {
			for k, room := range moves {
				if j+k >= len(finalMoves) {
					finalMoves = append(finalMoves, room+" ")
				} else {
					finalMoves[j+k] += room + " "
				}
			}
		}
	}
	return ProcessStrings(finalMoves)
}

func DeepCopyGraph(g *Graph) *Graph {
	newGraph := &Graph{Rooms: []*Room{}}
	for _, room := range g.Rooms {
		newGraph.Rooms = append(newGraph.Rooms, &Room{
			Roomname:    room.Roomname,
			Connections: make([]string, len(room.Connections)),
			Visited:     room.Visited,
		})
		copy(newGraph.Rooms[len(newGraph.Rooms)-1].Connections, room.Connections)
	}
	newGraph.StartRoomName = g.StartRoomName
	newGraph.EndRoomName = g.EndRoomName
	newGraph.Ants = g.Ants
	return newGraph
}

func LexicalLsort(s string) string {
	words := strings.Fields(s)
	sort.Slice(words, func(i, j int) bool {
		num1, _ := strconv.Atoi(strings.Split(words[i], "-")[0][1:])
		num2, _ := strconv.Atoi(strings.Split(words[j], "-")[0][1:])
		return num1 < num2
	})
	return strings.Join(words, " ")
}

func ProcessStrings(strs []string) []string {
	sortedStrs := make([]string, len(strs))
	for i, s := range strs {
		sortedStr := LexicalLsort(s)
		sortedStrs[i] = sortedStr
	}
	return sortedStrs
}

func contains(s []string, name string) bool {
	for _, str := range s {
		if str == name {
			return true
		}
	}
	return false
}

func shorterSearch(DFSSearch, BFSSearch []string) []string {
	if len(DFSSearch) > len(BFSSearch) {
		return BFSSearch
	}
	return DFSSearch
}

func validateFileGiveMeStrings() []string {
	// check if the file exists
	f, err := os.Stat(os.Args[1])
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("ERROR: invalid data format. File does not exist")
		}
	}
	// store the whole file in a string
	data, err := os.ReadFile(f.Name())
	if err != nil {
		log.Fatal("ERROR: invalid data format. File reading error", err)
	}
	// split the string into lines
	lines := strings.Split(string(data), "\n")
	// min lines need to be
	// 1. numofants
	// 2. startroom
	// 3. startroomname
	// 4. endroom
	// 5. endroomname
	// 6. onceconnection

	if len(lines) < 6 {
		log.Fatal("ERROR: invalid data format. Not enough lines")
	}

	// check if all the lines apart from the first one have a min length of 3
	for i := 1; i < len(lines); i++ {
		if len(lines[i]) < 3 {
			log.Fatal("ERROR: invalid data format. Line is too short")
		}
	}

	// check if the first line is a number
	numA, err := strconv.Atoi(lines[0])
	if err != nil {
		log.Fatal("ERROR: invalid data format. First line is not a number")
	}
	if numA <= 0 {
		log.Fatal("ERROR: invalid data format. Number of ants is negative or zero")
	}
	// check if there are more than two line with ##start or ##end
	startCount := 0
	endCount := 0
	for _, s := range lines {
		if strings.Contains(s, "##start") {
			startCount++
		}
		if strings.Contains(s, "##end") {
			endCount++
		}
	}
	if startCount > 1 || endCount > 1 {
		log.Fatal("ERROR: invalid data format. More than one ##start or ##end")
	}
	if startCount == 0 || endCount == 0 {
		log.Fatal("ERROR: invalid data format. No ##start or ##end")
	}
	// if there is any line which starts from #, and it's length is greater than 1 and it doesn't equal ##start or ##end, then remove that line from lines
	for i, s := range lines {
		if strings.HasPrefix(s, "#") && len(s) > 1 && !strings.Contains(s, "##start") && !strings.Contains(s, "##end") {
			lines = append(lines[:i], lines[i+1:]...)
		}
	}

	// check if there is atleast one line with ##start and one line with ##end
	for i, s := range lines {
		if strings.Contains(s, "##start") && i+1 < len(lines) {
			if !isValidRoomName(lines[i+1]) {
				log.Fatal("ERROR: invalid data format. Invalid room name")
			}
			break
		}
		if strings.Contains(s, "##end") && i+1 < len(lines) {
			if !isValidRoomName(lines[i+1]) {
				log.Fatal("ERROR: invalid data format. Invalid room name")
			}
			break
		}
	}
	chkConsRInTheEnd(lines)
	chkDuplicateCoords(lines)
	return lines
}

// check for connections in the end
func chkConsRInTheEnd(lines []string) {
	// check if the last line has a "-"
	for {
		if strings.Contains((lines[len(lines)-1]), "-") {
			lines = lines[:len(lines)-1]
		} else {
			break
		}
	}
	// now that we eliminated all the conitinous "-" contianing strings fromt he end, we are golden. Just now check if there is any other fucking line with "-"
	for _, s := range lines {
		if strings.Contains(s, "-") {
			log.Fatal("ERROR: invalid data format. Invalid connection, all connections have to be continuous in the end")
		}
	}
}

func chkDuplicateCoords(lines []string) {
	// check if there are duplicate coordinates
	coords := make(map[string]bool)
	countCoord := 0
	// get all the lines with coordinates
	for _, s := range lines {
		if strings.Contains(s, " ") {
			countCoord++
			words := strings.Fields(s)
			coords[words[1]+" "+words[2]] = true
		}
	}
	// check if the length of the map is equal to the number of lines with coordinates
	if len(coords) != countCoord {
		log.Fatal("ERROR: invalid data format. Duplicate coordinates")
	}
}

func isValidRoomName(name string) bool {
	// room is the the format name x y
	words := strings.Fields(name)
	if len(words) != 3 {
		return false
	}
	// check if the second and third word can be converted to int
	_, err := strconv.Atoi(words[1])
	if err != nil {
		return false
	}
	_, err = strconv.Atoi(words[2])
	return err == nil
}

//////////////////////////////////////////////////////////////////////////

//-FROOM

// checkUnconnectedRooms checks if there are rooms that are not connected to the anthill
func checkUnconnectedRooms(ah *AntHill) {
	for _, room := range ah.FRooms {
		if len(room.Connections) == 0 {
			NoGo(fmt.Sprintf("The room \"%v\" is not connected to the anthill", room.Name))
		}
	}
}

// AddConnections adds connections between rooms based on the given list of connections in incoming format ["room1-room2", "room2-room3", ...]
func AddConnections(OnlyConnections []string) {
	for _, connection := range OnlyConnections {
		room1Name := strings.Split(connection, "-")[0]
		room2Name := strings.Split(connection, "-")[1]
		room1 := GetRoomByName(room1Name)
		room2 := GetRoomByName(room2Name)
		room1.Connections = append(room1.Connections, room2)
		room2.Connections = append(room2.Connections, room1)
	}
}

func GetRoomByName(name string) *FRoom {
	for _, room := range ah.FRooms {
		if room.Name == name {
			return room
		}
	}
	return nil
}

// RemoveComments removes the comments from the original file lines
func RemoveComments(originalFileLines []string) []string {
	var filteredLines []string
	for _, line := range originalFileLines {
		if strings.HasPrefix(line, "#") && line != "##end" && line != "##start" {
			continue
		}
		filteredLines = append(filteredLines, line)
	}
	return filteredLines
}

// IsNumber checks if a string is a number
func IsNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// No2Dashes checks if there are 2 or more dashes in a line
func No2Dashes(s []string) {
	for _, line := range s {
		if len(strings.Split(line, "-")) > 2 {
			NoGo("2 or more dashes in a line are not allowed")
		}
	}
}

// No3Spaces checks if there are 3 or more spaces in a line
func No3Spaces(s []string) {
	for _, line := range s {
		if len(strings.Split(line, " ")) > 3 {
			NoGo("3 or more spaces in a line are not allowed")
		}
	}
}

// ExtractStartRoom extracts the start room from the slice
func ExtractStartRoom(s []string) {
	for i, line := range s {
		if line == "##start" {
			if i+1 < len(s) && IsRoom(s[i+1]) {
				ah.StartRoom = ConvertToRoom(s[i+1])
			} else {
				NoGo("")
			}
		}
	}
}

// ExtractEndRoom extracts the end room from the slice
func ExtractEndRoom(s []string) {
	for i, line := range s {
		if line == "##end" {
			if i+1 < len(s) && IsRoom(s[i+1]) {
				ah.EndRoom = ConvertToRoom(s[i+1])
			} else {
				NoGo("")
			}
		}
	}
}

// DeleteStartRoom deletes the start room from the slice
func DeleteStartRoom(s []string) []string {
	var filteredLines []string
	startRoomIndex := -1
	for i, line := range s {
		if i == startRoomIndex {
			continue
		}
		if line == "##start" {
			startRoomIndex = i + 1
			continue
		}
		filteredLines = append(filteredLines, line)
	}
	return filteredLines
}

// DeleteEndRoom deletes the end room from the slice
func DeleteEndRoom(s []string) []string {
	var filteredLines []string
	endRoomIndex := -1
	for i, line := range s {
		if i == endRoomIndex {
			continue
		}
		if line == "##end" {
			endRoomIndex = i + 1
			continue
		}
		filteredLines = append(filteredLines, line)
	}
	return filteredLines
}

// DeleteAllRooms deletes all the rooms from the slice
func DeleteAllRooms(s []string) []string {
	var filteredLines []string
	for _, line := range s {
		if !IsRoom(line) {
			filteredLines = append(filteredLines, line)
		}
	}
	return filteredLines
}

// NoDuplicateLines checks if there are duplicate lines in the slice
func NoDuplicateLines(s []string) {
	for i, line := range s {
		for j, line2 := range s {
			if i != j && line == line2 {
				NoGo("Duplicate lines are not allowed")
			}
		}
	}
}

// NoDuplicateCoordsOrNames checks if there are duplicate coordinates in the slice
func NoDuplicateCoordsOrNames(s []*FRoom) {
	for i, room := range s {
		for j, room2 := range s {
			if i != j && room.X == room2.X && room.Y == room2.Y {
				NoGo("Duplicate coordinates are not allowed")
			}
			if i != j && room.Name == room2.Name {
				NoGo("Duplicate room names are not allowed")
			}
		}
	}
}

// ExtractRooms extracts all the rooms from the slice
func ExtractRooms(s []string) {
	var rooms []*FRoom
	for _, line := range s {
		if IsRoom(line) {
			rooms = append(rooms, ConvertToRoom(line))
		}
	}
	rooms = append(rooms, ah.StartRoom)
	rooms = append(rooms, ah.EndRoom)
	NoDuplicateCoordsOrNames(rooms)
	// fill AntHill with the rooms
	ah.FRooms = rooms
}

// ConvertToRoom converts a string to a room
func ConvertToRoom(roomStr string) *FRoom {
	// split the room line into a slice
	roomStrSlice := strings.Split(roomStr, " ")
	// convert the coordinates to ints
	rName := roomStrSlice[0]
	x, _ := strconv.Atoi(roomStrSlice[1])
	y, _ := strconv.Atoi(roomStrSlice[2])
	return &FRoom{
		Name: rName,
		X:    x,
		Y:    y,
	}
}

// No # in last line, or it is a start or end room
func NoHashInLastLine(s []string) {
	if strings.HasPrefix(s[len(s)-1], "#") {
		NoGo("")
	}
}

// IsRoom checks if a string is a room
func IsRoom(s string) bool {
	return !((len(strings.Split(s, " ")) != 3) || !IsNumber(strings.Split(s, " ")[1]) || !IsNumber(strings.Split(s, " ")[2]))
}

func NoGo(msg string) {
	fmt.Println("ERROR: invalid data format")
	if msg != "" {
		fmt.Println("\033[101m" + msg + "\033[0m")
	}
	os.Exit(1)
}

func CheckRoomsInConnectionsPresent(OnlyConnections []string, AllRooms []string) {
	for _, connectionStr := range OnlyConnections {
		// split the connectionStr line into a slice of roomsnames by "-"
		roomNames := strings.Split(connectionStr, "-")
		if !Contains(AllRooms, roomNames[0]) || !Contains(AllRooms, roomNames[1]) {
			NoGo("ERROR: room in connection not present in rooms")
		}
	}
}

// Contains checks if a string is in a slice
func Contains(slice []string, elem string) bool {
	return strings.Contains(strings.Join(slice, "😎"), elem)
}

// CheckRoomsInConnectionsPresent checks if all the rooms in the connections are present in the rooms
func GetAllRoomNames(ah *AntHill) []string {
	var roomNames []string
	for _, room := range ah.FRooms {
		roomNames = append(roomNames, room.Name)
	}
	return roomNames
}

//-FROOM
