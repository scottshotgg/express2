// Includes:
// none

// Imports:
#include "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/defer.cpp"
#include <iostream>
#include <libgen.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <string>
#include <unistd.h>

// Namespaces:
namespace __file {
struct File {
  FILE *fp = nullptr;
  std::string mode = "";
  std::string name = "";
  std::string path = "";
  bool open = false;
  bool IsOpen() {
    defer onReturn, onExit;
    open = true;
    return open;
  }
  int Close() {
    defer onReturn, onExit;
    open = false;
    return fclose(fp);
  }
  int CurrentPostition() {
    defer onReturn, onExit;
    return ftell(fp);
  }
  int SeekPosition(int pos) {
    defer onReturn, onExit;
    return fseek(fp, pos, SEEK_SET);
  }
  int Length() {
    defer onReturn, onExit;
    onReturn.deferStack.push([&](...) { SeekPosition(CurrentPostition()); });
    fseek(fp, 0, SEEK_END);
    return CurrentPostition();
  }
  int Rename(std::string name) {
    defer onReturn, onExit;
    std::string newName = path + "/" + name;
    return Move(newName.c_str());
  }
  int Move(std::string name) {
    defer onReturn, onExit;
    std::string oldName = path + "/" + name;
    return rename(oldName.c_str(), name.c_str());
  }
  int Delete() {
    defer onReturn, onExit;
    std::string loc = path + "/" + name;
    return remove(loc.c_str());
  }
  std::string Read(int num) {
    defer onReturn, onExit;
    std::string line;
    int count = 0;
    char ch = fgetc(fp);
    while (!feof(fp) && count < num) {
      (count)++;
      line = line + ch;
      ch = fgetc(fp);
    }
    return line;
  }
  char ReadNextChar() {
    defer onReturn, onExit;
    if (!feof(fp)) {
      return fgetc(fp);
    }
    return '\0';
  }
  std::string ReadLine() {
    defer onReturn, onExit;
    std::string line;
    char ch = fgetc(fp);
    while (!feof(fp) && ch != '\n') {
      line = line + ch;
      ch = fgetc(fp);
    }
    return line;
  }
  int Write(std::string text) {
    defer onReturn, onExit;
    return fputs(text.c_str(), fp);
  }
  int WriteLine(std::string text) {
    defer onReturn, onExit;
    return Write(text + "\n");
  }
};
File Open(std::string loc, std::string mode) {
  defer onReturn, onExit;
  return File{
      .fp = fopen(loc.c_str(), mode.c_str()),
      .mode = mode,
      .name = loc,
      .path = loc,
      .open = true,
  };
}
std::string ReadFile(std::string loc) {
  defer onReturn, onExit;
  File f = Open(loc, "r");
  onReturn.deferStack.push([&](...) { f.Close(); });
  std::string contents;
  int amount = 100;
  char buff[amount];
  while (fgets(buff, amount, f.fp)) {
    contents = contents + std::string(buff);
  }
  return contents;
}
void WriteFile(std::string loc, std::string contents, bool overwrite) {
  defer onReturn, onExit;
  File f = Open(loc, "w");
  onReturn.deferStack.push([&](...) { f.Close(); });
  if (f.fp == nullptr || overwrite) {
    f.Write(contents);
  }
}
void CopyFile(std::string from, std::string to, bool overwrite) {
  defer onReturn, onExit;
  WriteFile(to, ReadFile(from), overwrite);
}
} // namespace __file

// Types:
// none

// Structs:

// Prototypes:
// none

// Functions:// none
// Main:
// generated: false
