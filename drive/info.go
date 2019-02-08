package drive

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"io"
	"encoding/json"
)

type FileInfoArgs struct {
	Out         io.Writer
	Id          string
	SizeInBytes bool
	JsonOutput  bool
}

func (self *Drive) Info(args FileInfoArgs) error {
	f, err := self.service.Files.Get(args.Id).Fields("id", "name", "size", "createdTime", "modifiedTime", "md5Checksum", "mimeType", "parents", "shared", "description", "webContentLink", "webViewLink").Do()
	if err != nil {
		return fmt.Errorf("Failed to get file: %s", err)
	}

	pathfinder := self.newPathfinder()
	absPath, err := pathfinder.absPath(f)
	if err != nil {
		return err
	}

	PrintFileInfo(PrintFileInfoArgs{
		Out:         args.Out,
		File:        f,
		Path:        absPath,
		SizeInBytes: args.SizeInBytes,
		JsonOutput:	 args.JsonOutput,
	})

	return nil
}

type PrintFileInfoArgs struct {
	Out         io.Writer
	File        *drive.File
	Path        string
	SizeInBytes bool
	JsonOutput	bool
}

func PrintFileInfo(args PrintFileInfoArgs) {
	f := args.File

	if args.JsonOutput {

		type FileInfo struct {
			Id            string
			Name          string
			Path          string
			Description   string
			Mime          string
			Size          string
			Created       string
			Modified      string
			Md5sum        string
			Shared        string
			Parents       string
			ViewUrl       string
			DownloadUrl   string
		}

		fileInfo := FileInfo {
			Id:            f.Id,
			Name:          f.Name,
			Path:          args.Path,
			Description:   f.Description,
			Mime:          f.MimeType,
			Size:          formatSize(f.Size, args.SizeInBytes),
			Created:       formatDatetime(f.CreatedTime),
			Modified:      formatDatetime(f.ModifiedTime),
			Md5sum:        f.Md5Checksum,
			Shared:        formatBool(f.Shared),
			Parents:       formatList(f.Parents),
			ViewUrl:       f.WebViewLink,
			DownloadUrl:   f.WebContentLink,
		}

		type JsonResult struct {
			Message   string
			Info      FileInfo
		}

		jsonResult := JsonResult {
			Message:   "Success",
			Info:      fileInfo,
		}

		b, err := json.MarshalIndent(jsonResult, "", "  ")
		if err != nil {
			jsonResult.Message = "Fail"
		}
		
		args.Out.Write(b)

	} else {

		items := []kv{
			kv{"Id", f.Id},
			kv{"Name", f.Name},
			kv{"Path", args.Path},
			kv{"Description", f.Description},
			kv{"Mime", f.MimeType},
			kv{"Size", formatSize(f.Size, args.SizeInBytes)},
			kv{"Created", formatDatetime(f.CreatedTime)},
			kv{"Modified", formatDatetime(f.ModifiedTime)},
			kv{"Md5sum", f.Md5Checksum},
			kv{"Shared", formatBool(f.Shared)},
			kv{"Parents", formatList(f.Parents)},
			kv{"ViewUrl", f.WebViewLink},
			kv{"DownloadUrl", f.WebContentLink},
		}

		for _, item := range items {
			if item.value != "" {
				fmt.Fprintf(args.Out, "%s: %s\n", item.key, item.value)
			}
		}
	}
}
