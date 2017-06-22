package main

import (
	"fmt"
	"log"
	"os"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"bytes"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/client"
	"io"
)

//structure for configuration
type SSHS3UploadAPI struct{
	Username string
	Host 	 string
	Pwd 	string
	Port 	string
	FilePath string
	SSHClient *ssh.Client
	SSHSession *ssh.Session
}

//Function to establish a connection from our local client
func(sshobject *SSHS3UploadAPI) dialConnection()(error){

	//Clientconfiguration
	sshConfig := &ssh.ClientConfig{
		User: sshobject.Username,
		Auth: []ssh.AuthMethod{ssh.Password(sshobject.Pwd)},
	}

	//Dialing a connection
	client, err := ssh.Dial("tcp", sshobject.Host, sshConfig)
	if err != nil {
		return  err
	}

	//creating a session
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return err
	}

	sshobject.SSHClient=client
	sshobject.SSHSession = session
	return err

}

//Function for S3 upload on a remote server
func(sshObject *SSHS3UploadAPI) uploadFileUsingSSHS3RemoteExecution()(error){

	//buffersize
	var buffSize= 128*128
	var fileBuffer bytes.Buffer
	var datafileSize int64

	//opening the file at our local server
	fileReference,err:=os.Open(sshObject.FilePath)

	sftpClient := sftp.NewClient(sshObject.SSHClient,sftp.MaxPacket(buffSize))

	flags := os.O_CREATE | os.O_WRONLY|os.O_TRUNC

	//opening the file at remote server
	remoteFilePath, fileErr := sftpClient.OpenFile("/data/REMOTEFILE.DAT",flags)

	//Copy the data to a remote file from our local server
	bytesRead,_:= io.Copy(remoteFilePath,io.LimitReader(fileReference,buffSize))

	for bytesRead !=0 && bytesRead > 0 {

		bytesRead,reader := io.Copy(remoteFilePath,io.LimitReader(fileReference,buffSize))

	}

	if err!=nil{

		fmt.Println("File not found in path",sshObject.FilePath)

		return err
	}

	fileStatsData,err2:=fileReference.Stat()

	if err2!=nil{

		fmt.Println("File stats not found in path",sshObject.FilePath)

		return err2
	}

	datafileSize = fileStatsData.Size()

	//Executing the command at the remote server for S3 upload via AWS CLI
	runCommand := fmt.Sprintf("aws s3 cp /data/REMOTETRANSACTION.DAT s3://danta/trasactions.DAT --expected-size %d",datafileSize)

	err = sshObject.SSHSession.Run(runCommand)
}

func main() {
	if len(os.Args) != 5 {
		log.Fatalf("Usage: %s  <host> <port> <user>  <PWD> <filePath>", os.Args[0])
	}

	//Read input arguments & create Ssh API object
	sshAPIObject:=SSHS3UploadAPI{}
	sshAPIObject.Host=os.Args[0]
	sshAPIObject.Port=os.Args[1]
	sshAPIObject.Username = os.Args[2]
	sshAPIObject.Pwd = os.Args[3]
	sshAPIObject.FilePath=os.Args[4]

	//Dialing a connection
	sshAPIObject.dialConnection()
	//File Upload
	sshAPIObject.uploadFileUsingSSHS3RemoteExecution()
}

