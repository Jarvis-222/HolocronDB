package server

import (
	"HolocronDB/config"
	"HolocronDB/core"
	"log"
	"net"
	"syscall"
)

func RunAsyncTCPServer() error {

	log.Println("Starting asynchronous T CP server...", config.Host, config.Port)

	max_clients := 20000
	con_clients := 0

	//Creating a slice to hold the events for Epoll
	var events []syscall.EpollEvent = make([]syscall.EpollEvent, max_clients)

	//,Creating the socket
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM|syscall.SOCK_NONBLOCK, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(serverFD)

	//Binding the socket to the address and port

	ip4 := net.ParseIP(config.Host).To4()

	if err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}); err != nil {
		return err
	}

	//Start listening
	if err = syscall.Listen(serverFD, max_clients); err != nil {
		return err
	}

	//Async IO starts here... !!!

	//Creating an Epoll instance
	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		return err
	}
	defer syscall.Close(epollFD)

	//Specifying the events we want to get hints about
	socketSeverEvent := syscall.EpollEvent{
		Events: syscall.EPOLLIN, //We want to know when the socket is ready for reading (i.e., when a new connection is incoming)
		Fd:     int32(serverFD), //The file descriptor of the server socket
	}

	//Registering the server socket with the Epoll instance
	if err = syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, serverFD, &socketSeverEvent); err != nil {
		return err
	}

	for {

		//See that if any FD is ready for an IO

		nevents, err := syscall.EpollWait(epollFD, events[:], -1)
		if err != nil {
			continue
		}

		for i := 0; i < nevents; i++ {

			if int(events[i].Fd) == serverFD {

				fd, _, err := syscall.Accept(serverFD)
				if err != nil {
					log.Println("Error accepting connection: ", err)
					continue
				}

				con_clients++

				socketClientEvent := syscall.EpollEvent{
					Events: syscall.EPOLLIN, //We want to know when the client socket is ready for reading (i.e., when a client sends a command)
					Fd:     int32(fd),       //The file descriptor of the client socket
				}

				if err = syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, fd, &socketClientEvent); err != nil {
					log.Println("Error adding client socket to epoll: ", err)
					syscall.Close(fd)
					continue
				}

				log.Println("New client connected with address: ", fd, " | con-current clients: ", con_clients)

			} else {

				comm := core.FDComm{Fd: int(events[i].Fd)}

				cmd, err := readCommand(comm)

				if err != nil {
					syscall.Close(int(events[i].Fd))
					log.Println("Client disconnected with address: ", events[i].Fd, " | con-current clients: ", con_clients-1)
					con_clients -= 1
					continue
				}
				respond(cmd, comm)

			}
		}

	}

}
