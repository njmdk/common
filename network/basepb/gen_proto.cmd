 for /r %%s in (*.proto) do (
 	protoc -I. -I%GOPATH%/src -I%GOPATH%/src/github.com/njmdk/common/pbbase --go_out . --proto_path . ./%%~ns.proto
 	protoc -I. -I%GOPATH%/src -I%GOPATH%/src/github.com/njmdk/common/pbbase --java_out %GOPATH%/src/github.com/njmdk/Server/pbmsg/java/ --proto_path . ./%%~ns.proto
 	protoc -I. -I%GOPATH%/src -I%GOPATH%/src/github.com/njmdk/common/pbbase --objc_out %GOPATH%/src/github.com/njmdk/Server/pbmsg/oc/basepb --proto_path . ./%%~ns.proto
 	protoc -I. -I%GOPATH%/src -I%GOPATH%/src/github.com/njmdk/common/pbbase --js_out %GOPATH%/src/github.com/njmdk/Server/pbmsg/js/basepb --proto_path . ./%%~ns.proto
 	protoc -I. -I%GOPATH%/src -I%GOPATH%/src/github.com/njmdk/common/pbbase --doc_out=. --doc_opt=json,%%~ns.json ./%%~ns.proto
 	CALL :CHECK_FAIL
 )

 :: ///
 :CHECK_FAIL
 @echo off
 if NOT ["%errorlevel%"]==["0"] (
     pause
     exit /b %errorlevel%
 )