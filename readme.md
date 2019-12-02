# Usage
trep

    -h  _print help_

    -d _directory to search_ ( default '.' )

    -e _extension to include_ ( default 'txt' )

    -r _the regular expression to search_ ( default '(.*?)' )

    -t _the regular expression to require_ ( default '.*' )

    -p _the replace. use go style sprintf syntax ( default %s[[0]] ) 

    -i _in match replacements_ replace characters in output matches

    
    
## Example

Assume a folder structure full of html files, and you want to list the sources of all image tags (assuming powershell)

`
trep -d . -e html -r 'src="""\\-/media/(.*?).ashx.*?"""' -p "%[1]s"
`

Assume a folder structure of text files with a csv of two columns, and you need to extract them to write sql.

`
trep -d . -e txt -r '(.*?) (.*?)' -p "insert into Table (Col1, Col2) values ('%[1]s', '%[2]s')" -i "':''" 
`

The in match replacement, since this is sql, replaces apostrophe with double apostrophe as to not allow the values to interfere with the sql execution.

The -d and -e parameters are given with the default values for illustration, but they won't be required if you're just doing the default values, obviously.