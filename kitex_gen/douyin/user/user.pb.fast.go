// Code generated by Fastpb v0.0.2. DO NOT EDIT.

package user

import (
	fmt "fmt"
	fastpb "github.com/cloudwego/fastpb"
)

var (
	_ = fmt.Errorf
	_ = fastpb.Skip
)

func (x *UserRequest) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
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
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_UserRequest[number], err)
}

func (x *UserRequest) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.UserId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *UserRequest) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.ActorId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *UserResponse) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
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
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_UserResponse[number], err)
}

func (x *UserResponse) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.StatusCode, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *UserResponse) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.StatusMsg, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *UserResponse) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	var v User
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.User = &v
	return offset, nil
}

func (x *User) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
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
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_User[number], err)
}

func (x *User) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.Id, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *User) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.Name, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *User) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	x.FollowCount, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *User) fastReadField4(buf []byte, _type int8) (offset int, err error) {
	x.FollowerCount, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *User) fastReadField5(buf []byte, _type int8) (offset int, err error) {
	x.IsFollow, offset, err = fastpb.ReadBool(buf, _type)
	return offset, err
}

func (x *UserRequest) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *UserRequest) fastWriteField1(buf []byte) (offset int) {
	if x.UserId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.UserId)
	return offset
}

func (x *UserRequest) fastWriteField2(buf []byte) (offset int) {
	if x.ActorId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 2, x.ActorId)
	return offset
}

func (x *UserResponse) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	return offset
}

func (x *UserResponse) fastWriteField1(buf []byte) (offset int) {
	if x.StatusCode == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.StatusCode)
	return offset
}

func (x *UserResponse) fastWriteField2(buf []byte) (offset int) {
	if x.StatusMsg == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.StatusMsg)
	return offset
}

func (x *UserResponse) fastWriteField3(buf []byte) (offset int) {
	if x.User == nil {
		return offset
	}
	offset += fastpb.WriteMessage(buf[offset:], 3, x.User)
	return offset
}

func (x *User) FastWrite(buf []byte) (offset int) {
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

func (x *User) fastWriteField1(buf []byte) (offset int) {
	if x.Id == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.Id)
	return offset
}

func (x *User) fastWriteField2(buf []byte) (offset int) {
	if x.Name == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.Name)
	return offset
}

func (x *User) fastWriteField3(buf []byte) (offset int) {
	if x.FollowCount == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 3, x.FollowCount)
	return offset
}

func (x *User) fastWriteField4(buf []byte) (offset int) {
	if x.FollowerCount == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 4, x.FollowerCount)
	return offset
}

func (x *User) fastWriteField5(buf []byte) (offset int) {
	if !x.IsFollow {
		return offset
	}
	offset += fastpb.WriteBool(buf[offset:], 5, x.IsFollow)
	return offset
}

func (x *UserRequest) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *UserRequest) sizeField1() (n int) {
	if x.UserId == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.UserId)
	return n
}

func (x *UserRequest) sizeField2() (n int) {
	if x.ActorId == 0 {
		return n
	}
	n += fastpb.SizeUint32(2, x.ActorId)
	return n
}

func (x *UserResponse) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	return n
}

func (x *UserResponse) sizeField1() (n int) {
	if x.StatusCode == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.StatusCode)
	return n
}

func (x *UserResponse) sizeField2() (n int) {
	if x.StatusMsg == "" {
		return n
	}
	n += fastpb.SizeString(2, x.StatusMsg)
	return n
}

func (x *UserResponse) sizeField3() (n int) {
	if x.User == nil {
		return n
	}
	n += fastpb.SizeMessage(3, x.User)
	return n
}

func (x *User) Size() (n int) {
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

func (x *User) sizeField1() (n int) {
	if x.Id == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.Id)
	return n
}

func (x *User) sizeField2() (n int) {
	if x.Name == "" {
		return n
	}
	n += fastpb.SizeString(2, x.Name)
	return n
}

func (x *User) sizeField3() (n int) {
	if x.FollowCount == 0 {
		return n
	}
	n += fastpb.SizeUint32(3, x.FollowCount)
	return n
}

func (x *User) sizeField4() (n int) {
	if x.FollowerCount == 0 {
		return n
	}
	n += fastpb.SizeUint32(4, x.FollowerCount)
	return n
}

func (x *User) sizeField5() (n int) {
	if !x.IsFollow {
		return n
	}
	n += fastpb.SizeBool(5, x.IsFollow)
	return n
}

var fieldIDToName_UserRequest = map[int32]string{
	1: "UserId",
	2: "ActorId",
}

var fieldIDToName_UserResponse = map[int32]string{
	1: "StatusCode",
	2: "StatusMsg",
	3: "User",
}

var fieldIDToName_User = map[int32]string{
	1: "Id",
	2: "Name",
	3: "FollowCount",
	4: "FollowerCount",
	5: "IsFollow",
}
