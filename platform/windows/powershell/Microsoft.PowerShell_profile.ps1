neofetch
# 0. preferences. some powershell settings here
# https://github.com/PowerShell/PSReadLine/blob/master/PSReadLine/SamplePSReadLineProfile.ps1 
# https://wentao.org/post/2022-07-10-psreadline/
# https://github.com/PowerShell/PSReadLine 
# https://zhuanlan.zhihu.com/p/137595941 
# https://asmcn.icopy.site/awesome/awesome-powershell/#awesome-powershell

# 1. Load PSReadLine and config it 
Import-Module PSReadLIne
Set-PSReadLineOption -PredictionSource History
Set-PSReadLineOption -HistorySearchCursorMovesToEnd
Set-PSReadLIneKeyHandler -Key "Tab" -Function MenuComplete
Set-PSReadLIneKeyHandler -Key "Ctrl+d" -Function ViExit
Set-PSReadLIneKeyHandler -Key "Ctrl+z" -Function Undo
Set-PSReadLIneKeyHandler -Key UpArrow -Function HistorySearchBackward
Set-PSReadLIneKeyHandler -Key DownArrow -Function HistorySearchForward

# 2. Load the plugin to optimize the ls function
# the DirColors can not exist with terminal-icons at the same time
# Import-Module DirColors
Import-Module Terminal-Icons

function ListDirectory {
	(Get-ChildItem).Name
	Write-Host("")
}
Set-Alias -Name ls -Value ListDirectory
Set-Alias -Name ll -Value Get-ChildItem

# 3. Load z-jumper to optimize the cd command
Import-Module Zlocation

# 4. Load the OMP to make cmd beautiful.
# you can use `Get-PoshThemes` to choose one theme you like.
Import-Module posh-git
oh-my-posh init pwsh | Invoke-Expression
oh-my-posh init pwsh --config 'C:\Users\Aiken\AppData\Local\Programs\oh-my-posh\themes\1_shell.omp.json' | Invoke-Expression

# 5.(to be finished) set alias for same experiment as linux
Set-Alias -name grep -value select-string
# Set-Alias -name grep -value findstr


# 6.  F7 to pop up a history page
Set-PSReadLineKeyHandler -Key F7 `
	-BriefDescription History `
	-LongDescription 'Show command history' `
	-ScriptBlock {
	$pattern = $null
	[Microsoft.PowerShell.PSConsoleReadLine]::GetBufferState([ref]$pattern, [ref]$null)
	if ($pattern) {
		$pattern = [regex]::Escape($pattern)
	}

	$history = [System.Collections.ArrayList]@(
		$last = ''
		$lines = ''
		foreach ($line in [System.IO.File]::ReadLines((Get-PSReadLineOption).HistorySavePath)) {
			if ($line.EndsWith('`')) {
				$line = $line.Substring(0, $line.Length - 1)
				$lines = if ($lines) {
					"$lines`n$line"
				}
				else {
					$line
				}
				continue
			}

			if ($lines) {
				$line = "$lines`n$line"
				$lines = ''
			}

			if (($line -cne $last) -and (!$pattern -or ($line -match $pattern))) {
				$last = $line
				$line
			}
		}
	)
	$history.Reverse()

	$command = $history | Out-GridView -Title History -PassThru
	if ($command) {
		[Microsoft.PowerShell.PSConsoleReadLine]::RevertLine()
		[Microsoft.PowerShell.PSConsoleReadLine]::Insert(($command -join "`n"))
	}
}

# 7. Network About Get-Key-Ipv4/Ipv6 Keyroute

function Get-IPv4Routes {
	Get-NetRoute -AddressFamily IPv4 | Where-Object -FilterScript {$_.NextHop -ne '0.0.0.0'}
}
Set-Alias -Name getrt -Value Get-IPv4Routes

function Get-IPv6Routes {
	Get-NetRoute -AddressFamily IPv6 | Where-Object -FilterScript {$_.NextHop -ne '::'}
}
Set-Alias -Name getrt6 -Value Get-IPv6Routes

function Get-IPv4Ipinfo {
    ipconfig /all | select-string WLAN -Context 0,18
}

function Get-IPv4ip{
    ipconfig /all | select-string "wlan" -Context 0,18 | Out-String -Stream |select-string Ipv4
}
Set-Alias -Name getipinfo -Value Get-IPv4Ipinfo
Set-Alias -Name getip -Value Get-IPv4ip

# some about select-string :https://blog.csdn.net/iteye_5455/article/details/127322684?spm=1001.2101.3001.6650.1&utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7EAD_ESQUERY%7Eyljh-1-127322684-blog-1795188.pc_relevant_3mothn_strategy_and_data_recovery&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7EAD_ESQUERY%7Eyljh-1-127322684-blog-1795188.pc_relevant_3mothn_strategy_and_data_recovery&utm_relevant_index=2 

# 8. set alias
# -- git
Set-Alias -Name g -Value git
# -- open dir
function OpenCurrentFolder {
	param
	(
		# 输入要打开的路径
		# 用法示例：open C:\
		# 默认路径：当前工作文件夹
		$Path = '.'
	)
	Invoke-Item $Path
}
Set-Alias -Name open -Value OpenCurrentFolder

# 9. set lvim 
Set-Alias lvim 'C:\Users\Aiken\.local\bin\lvim.ps1'

# 10. FIXME: make it params 
# set proxy for git and powershell  
function SetProxy4Git {
  git config --global http.proxy 127.0.0.1:8890
  git config --global https.proxy 127.0.0.1:8890
}
Set-Alias -Name gproxy -Value SetProxy4Git

function CancelProxy4Git {
  git config --global --unset http.proxy
  git config --global --unset https.proxy
}
Set-Alias -Name gcproxy -Value CancelProxy4Git

function SetProxy4PS {
  netsh winhttp set proxy 127.0.0.1:8890
}
Set-Alias -Name psproxy -Value SetProxy4PS 

function CancelProxy4PS {
  netsh winhttp reset proxy 127.0.0.1:8890
}
Set-Alias -Name pscproxy -Value CancelProxy4PS

