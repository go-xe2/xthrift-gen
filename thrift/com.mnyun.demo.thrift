namespace com.mnyun.demo
include 'com.mnyun.types'
typedef string str
typedef bool bl
typedef double dl
service helloService {
	helloResult sayHello(2:i32 age,1:str name);
}
