" Vim syntax file for Hydrogen programming language
" Language: Hydrogen
" Maintainer: Cyanrad
" Latest Revision: 2025

if exists("b:current_syntax")
  finish
endif

" Keywords
syn keyword hydrogenKeyword let fn print return
syn keyword hydrogenBuiltin filter map reduce len

" Operators
syn match hydrogenOperator "="
syn match hydrogenOperator "+"
syn match hydrogenOperator "++"
syn match hydrogenOperator "-"
syn match hydrogenOperator "--"
syn match hydrogenOperator "*"
syn match hydrogenOperator "/"
syn match hydrogenOperator "%"
syn match hydrogenOperator "=="
syn match hydrogenOperator "!="
syn match hydrogenOperator "<="
syn match hydrogenOperator ">="
syn match hydrogenOperator "<"
syn match hydrogenOperator ">"

" Delimiters
syn match hydrogenDelimiter "("
syn match hydrogenDelimiter ")"
syn match hydrogenDelimiter "\["
syn match hydrogenDelimiter "\]"
syn match hydrogenDelimiter "{"
syn match hydrogenDelimiter "}"
syn match hydrogenDelimiter ","
syn match hydrogenDelimiter ";"

" Numbers
syn match hydrogenNumber "\<\d\+\>"
syn match hydrogenNumber "\<\d\+\.\d*\>"

" Strings
syn region hydrogenString start='"' end='"' contained

" Comments (if your language supports them)
syn match hydrogenComment "#.*$"

" Identifiers (variable names, function names)
syn match hydrogenIdentifier "\<[a-zA-Z_][a-zA-Z0-9_]*\>"

" Function calls
syn match hydrogenFunction "\<[a-zA-Z_][a-zA-Z0-9_]*\>\s*("he=e-1

" Boolean values
syn keyword hydrogenBoolean true false

" Define highlighting groups
hi def link hydrogenKeyword     Keyword
hi def link hydrogenBuiltin     Function
hi def link hydrogenOperator    Operator
hi def link hydrogenDelimiter   Delimiter
hi def link hydrogenNumber      Number
hi def link hydrogenString      String
hi def link hydrogenComment     Comment
hi def link hydrogenIdentifier  Identifier
hi def link hydrogenFunction    Function
hi def link hydrogenProperty    Special
hi def link hydrogenBoolean     Boolean

let b:current_syntax = "hydrogen"

" should be in a seperae file for filetype detection in .config/nvim/ftdetect/hydrogen.vim
au BufRead,BufNewFile *.hy setfiletype hydrogen
