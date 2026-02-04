@echo off
setlocal

if not defined workspaceFolder (
    if defined WORKSPACE_FOLDER (
        set "workspaceFolder=%WORKSPACE_FOLDER%"
    ) else (
        for %%I in ("%~dp0..\..") do set "workspaceFolder=%%~fI"
    )
)

echo workspaceFolder: %workspaceFolder%

if not defined workspaceFolder (
    echo workspaceFolder is empty. skip delete.
    exit /b 0
)

del /f /q "%workspaceFolder%\cmd\__debug_bin*" >nul 2>&1
del /f /q "%workspaceFolder%\go\cmd\__debug_bin*" >nul 2>&1

exit /b 0
