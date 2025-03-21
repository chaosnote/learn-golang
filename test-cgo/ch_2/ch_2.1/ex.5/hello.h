// hello.h
// 为了适配 CGO 导出的 C 语言函数，我们禁止了在函数的声明语句中的 const 修饰符
void SayHello(/*const*/ char* s);
