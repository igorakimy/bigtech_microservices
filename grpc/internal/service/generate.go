package service

//go:generate sh -c "rm -rf mocks && mkdir mocks"
//go:generate ../../bin/minimock -i NoteService -o ./mocks/ -s "_minimock.go"
