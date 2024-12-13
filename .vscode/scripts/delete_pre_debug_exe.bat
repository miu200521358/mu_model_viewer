echo off

echo workspaceFolder: %workspaceFolder%

call del %workspaceFolder%\cmd\__debug_bin*
call del %workspaceFolder%\go\cmd\__debug_bin*

exit 0
