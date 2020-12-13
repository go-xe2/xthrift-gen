namespace com.mnyun.types
typedef string str
typedef bool bl
typedef double dl
typedef list<map<str,str>> rows
typedef i8 int
struct helloData {
	1:optional str name;
	2:bl sex;
}
struct helloResult {
	1:i32 status;
	2:str msg;
	3:helloData data;
}
