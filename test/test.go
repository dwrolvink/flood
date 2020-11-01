package test

type TestObj struct {
	Text string
	Amount int
}

func NewTestObj(text string) *TestObj {  
	var obj = TestObj{Text: text}
	obj.Amount = 15

	return &obj
}