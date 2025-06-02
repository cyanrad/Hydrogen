# Hydrogen
![Hydrogen](https://raw.githubusercontent.com/cyanrad/Hydrogen/refs/heads/master/common/src/images/hydrogen-dark.png "Hydrogen")
A simple python-like toy programming language that takes inspiration from functional languages.

## Examples
This is an example that showcases the basic features of the language. 

### REPL
You can use the basic REPL to evaluate expressions and see results immediately.
```bash
git clone https://github.com/cyanrad/Hydrogen
cd Hydrogen
go run .
>> print("Hello, World!")
"Hello, World!"
```

### Library Program
you can run Hydrogen programs from filse using the `-file` flag.
```bash
go run . -file library.hy
```

The following is a simple library management toy program written in Hydrogen.

```js
# Basic Library Management System - Showcasing functions, arrays, hashes, and functional style

let books = [
  {"id": 1, "title": "1984", "author": "Orwell", "available": true, "pages": 328},
  {"id": 2, "title": "Dune", "author": "Herbert", "available": false, "pages": 688},
  {"id": 3, "title": "Neuromancer", "author": "Gibson", "available": true, "pages": 271},
  {"id": 4, "title": "Foundation", "author": "Asimov", "available": true, "pages": 244}
];

let get_available_books = fn (book_list) {
    filter(book_list, fn (book) { book["available"]; });
};

let create_reading_list = fn (book_list) {
    let available = get_available_books(book_list);
    return map(available, fn (book) {book["title"]});
};

let total_pages = fn (book_list) {
    reduce(book_list, 0, fn (acc, book) {acc + book["pages"]});
};

let get_by_author = fn (book_list, author) {
    filter(book_list, fn (book) { book["author"] == author; });
};

let available_books = get_available_books(books);
let reading_list = create_reading_list(books);
let total_book_pages = total_pages(books);

print("Available books:");
print(available_books);

print("Reading list:");
print(reading_list);

print("Total pages in library:");
print(total_book_pages);

print("Books by Orwell:");
print(get_by_author(books, "Orwell"));
```
