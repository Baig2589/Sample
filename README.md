# Sample Task - Establishing a connection to a remote server via SSH, copy a file from our local server to a remote destination and upload that copied file to a S3 bucket.
Assumptions - The remote server is configured to a AWS IAM profile.

-> Here we made use of two packages namely "github.com/pkg/sftp" and "golang.org/x/crypto/ssh" for establishing a connection to a remote        serever and to file transfer.

-> We passed <host> - Hostname, <port> - Port No, <user> - Username,  <PWD> - Password and  <filePath> -  Filepath to establish a connection using "dialConnection" function. To capture SSHSession and SSHClient Objects.

-> We called uploadFileUsingSSHS3RemoteExecution for file Transfer to a remote server and finally for a S3 Upload.


