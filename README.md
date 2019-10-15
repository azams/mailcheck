Account verifier
Usage:
	-s		- Define Custom SMTP server.
			    If the value is empty, it will check based on email domain.
			    Support: gmail,yahoo,aol,hotmail,icloud,outlook
	-x		- Checking for specific domain (separated by comma)
			    Support: gmail,yahoo,aol,hotmail,icloud,outlook
	-f		- Define list of email password to check.
	-p		- Define SMTP port (default: 587).
	-d		- Define delimiter for email & password. (default is ':')'
	-m		- Mode can be 'smtp' or 'wordpress'. (default: smtp)
			    When using wordpress mode, server must be a url with scheme.

Example 1 : ./check -s smtp.example.com -p 587 -f lists.txt -d '|'
Example 2 : ./check -f lists.txt -d '|'
Example 3 : ./check -f lists.txt -d '|' -x aol,gmail,icloud
Example 4 : ./check -f lists.txt -d '|' -s https://target.com -m wordpress
