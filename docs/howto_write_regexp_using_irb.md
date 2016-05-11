# Howto write regexp using IRB

When you construct your regexp, the Ruby IRB is helpful. To start out with you can comment out all task systems except for the `none` system to verify that
you capture the correct task id's. Use the interactive ruby interpreter to test your regular expressions.

Example one, checking the above example regexp works as we expect

	2.2.1 :001 > commit_message ="Commit message header
	2.2.1 :002"> 
	2.2.1 :003"> More lines in commit message
	2.2.1 :004"> Issue: 12"
	 => "Commit message header\n\nMore lines in commit message\nIssue: 12"
 	2.2.1 :015 > regex = /Issue:\s+(\d+)/im
 	=> /Issue:\s+(\d+)/mi 
	2.2.1 :016 > commit_message.scan(regex)
 	=> [["12"]] 

The script uses the scan method to scan the commit mesage, so for each result the script looks in the first capture group, in this case
PAC would have picked up an issue with the id "12". The scan method can return multiple matches on the same text so make sure your capture first capture group
contains the id you want.



