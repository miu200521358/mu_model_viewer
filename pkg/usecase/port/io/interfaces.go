// 指示: miu200521358
package io

import portio "github.com/miu200521358/mlib_go/pkg/usecase/port/io"

// IFileReader は入出力共通の読み込み契約を表す。
type IFileReader = portio.IFileReader

// IFileWriter は入出力共通の書き込み契約を表す。
type IFileWriter = portio.IFileWriter

// IFileRepository は読み書きの共通契約を表す。
type IFileRepository = portio.IFileRepository

// IPathService は出力パス生成の契約を表す。
type IPathService = portio.IPathService

// ITextureValidator はテクスチャ検証の契約を表す。
type ITextureValidator = portio.ITextureValidator

// SaveOptions は保存時のオプションを表す。
type SaveOptions = portio.SaveOptions
