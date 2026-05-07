" reference the web and the tencent get the best setting and use it always
" suppose the mouse operation
set mouse=a
set selection=exclusive
set selectmode=mouse,key

" set the clipboard
set clipboard+=unnamed

" set the ruler line of vim
set ruler
set magic
set showcmd
set ttyfast
set wildmenu
set showmode
set lazyredraw
set cmdheight=2
set laststatus=2
set completeopt=preview,menu
set backspace=indent,eol,start

" detect the file type
filetype on
filetype plugin indent on

" set the theme
set t_Co=256
autocmd vimenter * hi Normal guibg=NONE ctermbg=NONE " transparent bg

" set the basic function
set nu
set nowrap
set confirm
set history=1000
set timeoutlen=350

" set the scoll rule
set scrolljump=5
set scrolloff=1

" the indent
set ai
"set si
"set cindent
autocmd Filetype cpp setlocal si cindent

" the tab setting
set smarttab
set tabstop=4
set shiftwidth=4
set softtabstop=4

autocmd Filetype cpp setlocal expandtab tabstop=2 shiftwidth=2
autocmd Filetype python setlocal expandtab tabstop=4 shiftwidth=4

" the search function and hilight
set showmatch " show match brackets
set smartcase
set ignorecase
set nohlsearch
set incsearch

" avoid compatible problem of vim and vi
set nocompatible

" avoid the tmp file and sup file
set noswapfile
set nobackup

" highlight the syntax
syntax on


" suppost chinese encoding
set encoding=utf-8
set termencoding=utf-8
set fileencodings=utf-8,ucs-bom,gb18030,gbk,gb2312,cp936
set list lcs=trail:·,tab:»·,nbsp:.,extends:#

" auto load when file change
set autoread

" reference : https://blog.csdn.net/strategycn/article/details/7620261
set statusline=%F%m%r%h%w\ [FORMAT=%{&ff}]\ [TYPE=%Y]\ [POS=%04l,%04v][%p%%]\ [LEN=%L]\ [TIME=%{strftime('%c')}]
" :set statusline=%2*%n%m%r%h%w%*\ %F\ %1*[FORMAT=%2*%{&ff}:%{&fenc!=''?&fenc:&enc}%1*]\ [TYPE=%2*%Y%1*]\ [COL=%2*%03v%1*]\ [ROW=%2*%03l%1*/%3*%L(%p%%)%1*]\

" [KEYMAPPING PART]
" set leaderkey to space
let mapleader="\<space>"

" set tab and shift to change indent
nmap <tab> V>
nmap <S-tab> V<
vmap <tab> >gv
vmap <S-tab> <gv

" del the space in end.
nnoremap <leader>de :%s/\s\+$//<cr>:let @/=''<CR>
" replace space as tab to solve the mix indent problem
nnoremap <leader>ds :retab!<CR>
" edit the nvim config file
nnoremap <leader>ev :vsp $MYVIMRC<CR>
" show diff <filename>
nnoremap <leader>df :vert diffsplit
" list file tree
map <leader>ft :tabnew .<cr>
" close tab
nmap <leader>tq :bp<cr>:bd #<cr>

" switch between windows
nmap <silent> <C>k :wincmd k<CR>
nmap <silent> <C>j :wincmd j<CR>
nmap <silent> <C>h :wincmd h<CR>
nmap <silent> <C>l :wincmd l<CR>

" basic keymapping
nnoremap <leader>q :q<CR>
nnoremap <leader>qa :qa<CR>
nnoremap <leader>w :w<CR>
nnoremap <leader>wa :wa<CR>

" tab operation
nmap <leader>ts :tabs<cr>
nmap <leader>tq :tabclose<cr>
nmap <leader>tn :tabn<cr>
nmap <leader>tp :tabp<cr>

" add some keyshot in command mode
nmap <C-a> <Home>
nmap <C-e> <End>
nmap <C-p> <PageUp>
nmap <C-n> <PageDown>

" [AUTOCMD PART]
" compile and run script
map <F5> :call CompileRunGcc()<CR>
func! CompileRunGcc()
	exec "w"
	if &filetype == 'cpp'
		exec '!g++ % -o %<'
		exec '!time ./%<'
		exec '!rm ./%<'
	elseif &filetype == 'python'
		exec '!python %'
	elseif &filetype == 'sh'
		:!time sh %
	endif
endfunc

" add header for py,cpp,sh
autocmd BufNewFile *.cpp,*.[ch],*.sh,*.py,*.md exec ":call SetTitle()"
func SetTitle()
	"如果文件类型为.sh文件
	if &filetype == 'sh'
		call setline(1, "# File Name: ".expand("%"))
		call append(line("."), "# Author: AikenHong")
		call append(line(".")+1, "# mail: h.aiken.970@gmail.com")
		call append(line(".")+2, "# Created Time: ".strftime("%c"))
		call append(line(".")+3, "")
	endif
	if &filetype == 'cpp'
		call setline(1, "/*")
		call append(line("."), "# File Name: ".expand("%"))
		call append(line(".")+1, "# Author: AikenHong")
		call append(line(".")+2, "# mail: h.aiken.970@gmail.com")
		call append(line(".")+3, "# Created Time: ".strftime("%c"))
		call append(line(".")+4, " */")
		call append(line(".")+5, " ")
		call append(line(".")+6, "#include <iostream>")
		call append(line(".")+7, "#include <algorithm>")
		call append(line(".")+8, "#include <vector>")
		call append(line(".")+9, "#include <stack>")
		call append(line(".")+10, "#include <queue>")
		call append(line(".")+11, "#include <list>")
		call append(line(".")+12, "#include <map>")
		call append(line(".")+13, "#include <cmath>")
		call append(line(".")+14, "#include <set>")
		call append(line(".")+15, "")
		call append(line(".")+16, "using namespace std;")
		call append(line(".")+17, "")
		call append(line(".")+18, "int main()")
		call append(line(".")+19, "{")
		call append(line(".")+20, "    ")
		call append(line(".")+21, "    ")
		call append(line(".")+22, "    return 0;")
		call append(line(".")+23, "}")
	endif
	if &filetype == 'python'
		call setline(1, "\"\"\"")
		call append(line("."), "# File Name: ".expand("%"))
		call append(line(".")+1, "# Author: AikenHong")
		call append(line(".")+2, "# mail: h.aiken.970@gmail.com")
		call append(line(".")+3, "# Created Time: ".strftime("%c"))
		call append(line(".")+4, "\"\"\"")
		call append(line(".")+5, "")
	endif
	if &filetype == 'markdown'
		call setline(1,"---")
		call append(line("."), "title: ")
		call append(line(".")+1, "subtitle: ")
		call append(line(".")+2, "toc: true")
		call append(line(".")+3, "lang: cn ")
		call append(line(".")+4, "catalog: true")
		call append(line(".")+5, "date: ".strftime("%Y-%m-%d %H:%M:%s"))
		call append(line(".")+6, "cover: /img/header_img/lml_bg.jpg")
		call append(line(".")+7, "thumbnail: /img/header_img/lml_bg.jpg")
		call append(line(".")+8, "mathjax: false")
		call append(line(".")+9, "tag: ")
		call append(line(".")+10, "categories: ")
		call append(line(".")+11, "---")
	endif
	"新建文件后，自动定位到文件末尾
	autocmd BufNewFile * normal G
endfunction
