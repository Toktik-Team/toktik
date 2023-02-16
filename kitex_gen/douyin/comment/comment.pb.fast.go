// Code generated by Fastpb v0.0.2. DO NOT EDIT.

package comment

import (
	fmt "fmt"
	fastpb "github.com/cloudwego/fastpb"
	user "toktik/kitex_gen/douyin/user"
)

var (
	_ = fmt.Errorf
	_ = fastpb.Skip
)

func (x *Comment) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 3:
		offset, err = x.fastReadField3(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 4:
		offset, err = x.fastReadField4(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_Comment[number], err)
}

func (x *Comment) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.Id, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *Comment) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	var v user.User
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.User = &v
	return offset, nil
}

func (x *Comment) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	x.Content, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *Comment) fastReadField4(buf []byte, _type int8) (offset int, err error) {
	x.CreateDate, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *ActionCommentRequest) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 3:
		offset, err = x.fastReadField3(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 4:
		offset, err = x.fastReadField4(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 5:
		offset, err = x.fastReadField5(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_ActionCommentRequest[number], err)
}

func (x *ActionCommentRequest) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.ActorId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *ActionCommentRequest) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.VideoId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *ActionCommentRequest) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	var v int32
	v, offset, err = fastpb.ReadInt32(buf, _type)
	if err != nil {
		return offset, err
	}
	x.ActionType = ActionCommentType(v)
	return offset, nil
}

func (x *ActionCommentRequest) fastReadField4(buf []byte, _type int8) (offset int, err error) {
	var ov ActionCommentRequest_CommentText
	x.Action = &ov
	ov.CommentText, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *ActionCommentRequest) fastReadField5(buf []byte, _type int8) (offset int, err error) {
	var ov ActionCommentRequest_CommentId
	x.Action = &ov
	ov.CommentId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *ActionCommentResponse) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 3:
		offset, err = x.fastReadField3(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_ActionCommentResponse[number], err)
}

func (x *ActionCommentResponse) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.StatusCode, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *ActionCommentResponse) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	tmp, offset, err := fastpb.ReadString(buf, _type)
	x.StatusMsg = &tmp
	return offset, err
}

func (x *ActionCommentResponse) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	var v Comment
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.Comment = &v
	return offset, nil
}

func (x *ListCommentRequest) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_ListCommentRequest[number], err)
}

func (x *ListCommentRequest) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.ActorId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *ListCommentRequest) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.VideoId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *ListCommentResponse) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 3:
		offset, err = x.fastReadField3(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_ListCommentResponse[number], err)
}

func (x *ListCommentResponse) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.StatusCode, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *ListCommentResponse) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	tmp, offset, err := fastpb.ReadString(buf, _type)
	x.StatusMsg = &tmp
	return offset, err
}

func (x *ListCommentResponse) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	var v Comment
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.CommentList = append(x.CommentList, &v)
	return offset, nil
}

func (x *CountCommentRequest) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_CountCommentRequest[number], err)
}

func (x *CountCommentRequest) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.ActorId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *CountCommentRequest) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.VideoId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *CountCommentResponse) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 3:
		offset, err = x.fastReadField3(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_CountCommentResponse[number], err)
}

func (x *CountCommentResponse) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.StatusCode, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *CountCommentResponse) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	tmp, offset, err := fastpb.ReadString(buf, _type)
	x.StatusMsg = &tmp
	return offset, err
}

func (x *CountCommentResponse) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	x.CommentCount, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *Comment) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	offset += x.fastWriteField4(buf[offset:])
	return offset
}

func (x *Comment) fastWriteField1(buf []byte) (offset int) {
	if x.Id == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.Id)
	return offset
}

func (x *Comment) fastWriteField2(buf []byte) (offset int) {
	if x.User == nil {
		return offset
	}
	offset += fastpb.WriteMessage(buf[offset:], 2, x.User)
	return offset
}

func (x *Comment) fastWriteField3(buf []byte) (offset int) {
	if x.Content == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 3, x.Content)
	return offset
}

func (x *Comment) fastWriteField4(buf []byte) (offset int) {
	if x.CreateDate == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 4, x.CreateDate)
	return offset
}

func (x *ActionCommentRequest) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	offset += x.fastWriteField4(buf[offset:])
	offset += x.fastWriteField5(buf[offset:])
	return offset
}

func (x *ActionCommentRequest) fastWriteField1(buf []byte) (offset int) {
	if x.ActorId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.ActorId)
	return offset
}

func (x *ActionCommentRequest) fastWriteField2(buf []byte) (offset int) {
	if x.VideoId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 2, x.VideoId)
	return offset
}

func (x *ActionCommentRequest) fastWriteField3(buf []byte) (offset int) {
	if x.ActionType == 0 {
		return offset
	}
	offset += fastpb.WriteInt32(buf[offset:], 3, int32(x.ActionType))
	return offset
}

func (x *ActionCommentRequest) fastWriteField4(buf []byte) (offset int) {
	if x.GetCommentText() == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 4, x.GetCommentText())
	return offset
}

func (x *ActionCommentRequest) fastWriteField5(buf []byte) (offset int) {
	if x.GetCommentId() == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 5, x.GetCommentId())
	return offset
}

func (x *ActionCommentResponse) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	return offset
}

func (x *ActionCommentResponse) fastWriteField1(buf []byte) (offset int) {
	if x.StatusCode == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.StatusCode)
	return offset
}

func (x *ActionCommentResponse) fastWriteField2(buf []byte) (offset int) {
	if x.StatusMsg == nil {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, *x.StatusMsg)
	return offset
}

func (x *ActionCommentResponse) fastWriteField3(buf []byte) (offset int) {
	if x.Comment == nil {
		return offset
	}
	offset += fastpb.WriteMessage(buf[offset:], 3, x.Comment)
	return offset
}

func (x *ListCommentRequest) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *ListCommentRequest) fastWriteField1(buf []byte) (offset int) {
	if x.ActorId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.ActorId)
	return offset
}

func (x *ListCommentRequest) fastWriteField2(buf []byte) (offset int) {
	if x.VideoId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 2, x.VideoId)
	return offset
}

func (x *ListCommentResponse) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	return offset
}

func (x *ListCommentResponse) fastWriteField1(buf []byte) (offset int) {
	if x.StatusCode == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.StatusCode)
	return offset
}

func (x *ListCommentResponse) fastWriteField2(buf []byte) (offset int) {
	if x.StatusMsg == nil {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, *x.StatusMsg)
	return offset
}

func (x *ListCommentResponse) fastWriteField3(buf []byte) (offset int) {
	if x.CommentList == nil {
		return offset
	}
	for i := range x.CommentList {
		offset += fastpb.WriteMessage(buf[offset:], 3, x.CommentList[i])
	}
	return offset
}

func (x *CountCommentRequest) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *CountCommentRequest) fastWriteField1(buf []byte) (offset int) {
	if x.ActorId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.ActorId)
	return offset
}

func (x *CountCommentRequest) fastWriteField2(buf []byte) (offset int) {
	if x.VideoId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 2, x.VideoId)
	return offset
}

func (x *CountCommentResponse) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	return offset
}

func (x *CountCommentResponse) fastWriteField1(buf []byte) (offset int) {
	if x.StatusCode == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.StatusCode)
	return offset
}

func (x *CountCommentResponse) fastWriteField2(buf []byte) (offset int) {
	if x.StatusMsg == nil {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, *x.StatusMsg)
	return offset
}

func (x *CountCommentResponse) fastWriteField3(buf []byte) (offset int) {
	if x.CommentCount == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 3, x.CommentCount)
	return offset
}

func (x *Comment) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	n += x.sizeField4()
	return n
}

func (x *Comment) sizeField1() (n int) {
	if x.Id == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.Id)
	return n
}

func (x *Comment) sizeField2() (n int) {
	if x.User == nil {
		return n
	}
	n += fastpb.SizeMessage(2, x.User)
	return n
}

func (x *Comment) sizeField3() (n int) {
	if x.Content == "" {
		return n
	}
	n += fastpb.SizeString(3, x.Content)
	return n
}

func (x *Comment) sizeField4() (n int) {
	if x.CreateDate == "" {
		return n
	}
	n += fastpb.SizeString(4, x.CreateDate)
	return n
}

func (x *ActionCommentRequest) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	n += x.sizeField4()
	n += x.sizeField5()
	return n
}

func (x *ActionCommentRequest) sizeField1() (n int) {
	if x.ActorId == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.ActorId)
	return n
}

func (x *ActionCommentRequest) sizeField2() (n int) {
	if x.VideoId == 0 {
		return n
	}
	n += fastpb.SizeUint32(2, x.VideoId)
	return n
}

func (x *ActionCommentRequest) sizeField3() (n int) {
	if x.ActionType == 0 {
		return n
	}
	n += fastpb.SizeInt32(3, int32(x.ActionType))
	return n
}

func (x *ActionCommentRequest) sizeField4() (n int) {
	if x.GetCommentText() == "" {
		return n
	}
	n += fastpb.SizeString(4, x.GetCommentText())
	return n
}

func (x *ActionCommentRequest) sizeField5() (n int) {
	if x.GetCommentId() == 0 {
		return n
	}
	n += fastpb.SizeUint32(5, x.GetCommentId())
	return n
}

func (x *ActionCommentResponse) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	return n
}

func (x *ActionCommentResponse) sizeField1() (n int) {
	if x.StatusCode == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.StatusCode)
	return n
}

func (x *ActionCommentResponse) sizeField2() (n int) {
	if x.StatusMsg == nil {
		return n
	}
	n += fastpb.SizeString(2, *x.StatusMsg)
	return n
}

func (x *ActionCommentResponse) sizeField3() (n int) {
	if x.Comment == nil {
		return n
	}
	n += fastpb.SizeMessage(3, x.Comment)
	return n
}

func (x *ListCommentRequest) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *ListCommentRequest) sizeField1() (n int) {
	if x.ActorId == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.ActorId)
	return n
}

func (x *ListCommentRequest) sizeField2() (n int) {
	if x.VideoId == 0 {
		return n
	}
	n += fastpb.SizeUint32(2, x.VideoId)
	return n
}

func (x *ListCommentResponse) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	return n
}

func (x *ListCommentResponse) sizeField1() (n int) {
	if x.StatusCode == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.StatusCode)
	return n
}

func (x *ListCommentResponse) sizeField2() (n int) {
	if x.StatusMsg == nil {
		return n
	}
	n += fastpb.SizeString(2, *x.StatusMsg)
	return n
}

func (x *ListCommentResponse) sizeField3() (n int) {
	if x.CommentList == nil {
		return n
	}
	for i := range x.CommentList {
		n += fastpb.SizeMessage(3, x.CommentList[i])
	}
	return n
}

func (x *CountCommentRequest) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *CountCommentRequest) sizeField1() (n int) {
	if x.ActorId == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.ActorId)
	return n
}

func (x *CountCommentRequest) sizeField2() (n int) {
	if x.VideoId == 0 {
		return n
	}
	n += fastpb.SizeUint32(2, x.VideoId)
	return n
}

func (x *CountCommentResponse) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	return n
}

func (x *CountCommentResponse) sizeField1() (n int) {
	if x.StatusCode == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.StatusCode)
	return n
}

func (x *CountCommentResponse) sizeField2() (n int) {
	if x.StatusMsg == nil {
		return n
	}
	n += fastpb.SizeString(2, *x.StatusMsg)
	return n
}

func (x *CountCommentResponse) sizeField3() (n int) {
	if x.CommentCount == 0 {
		return n
	}
	n += fastpb.SizeUint32(3, x.CommentCount)
	return n
}

var fieldIDToName_Comment = map[int32]string{
	1: "Id",
	2: "User",
	3: "Content",
	4: "CreateDate",
}

var fieldIDToName_ActionCommentRequest = map[int32]string{
	1: "ActorId",
	2: "VideoId",
	3: "ActionType",
	4: "CommentText",
	5: "CommentId",
}

var fieldIDToName_ActionCommentResponse = map[int32]string{
	1: "StatusCode",
	2: "StatusMsg",
	3: "Comment",
}

var fieldIDToName_ListCommentRequest = map[int32]string{
	1: "ActorId",
	2: "VideoId",
}

var fieldIDToName_ListCommentResponse = map[int32]string{
	1: "StatusCode",
	2: "StatusMsg",
	3: "CommentList",
}

var fieldIDToName_CountCommentRequest = map[int32]string{
	1: "ActorId",
	2: "VideoId",
}

var fieldIDToName_CountCommentResponse = map[int32]string{
	1: "StatusCode",
	2: "StatusMsg",
	3: "CommentCount",
}

var _ = user.File_user_proto
