# Function to load environment variables from .envrc file
function Load-DirenvVars {
    param (
        [string]$envrcFile
    )

    # Check if the .envrc file exists
    if (-Not (Test-Path $envrcFile)) {
        Write-Host "Error: .envrc file not found at $envrcFile"
        return
    }

    # Read the content of the .envrc file
    $lines = Get-Content $envrcFile

    foreach ($line in $lines) {
        # Ignore empty lines and comments
        if ($line.Trim() -match '^\s*(#.*|$)') {
            continue
        }

        # Look for export statements like: export VAR_NAME=value
        if ($line -match '^\s*export\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*=\s*(.*)\s*$') {
            $varName = $matches[1]
            $varValue = $matches[2]

            # Set the environment variable in the PowerShell session
            [System.Environment]::SetEnvironmentVariable($varName, $varValue, [System.EnvironmentVariableTarget]::Process)

            # Optional: Output the set environment variable for confirmation
            Write-Host "Set environment variable: $varName = $varValue"
        }
    }
}

# Usage example: Load variables from .envrc
$envrcFile = ".\.envrc" # Change this to the path of your .envrc file
Load-DirenvVars -envrcFile $envrcFile
