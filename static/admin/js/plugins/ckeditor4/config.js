/**
 * @license Copyright (c) 2003-2019, CKSource - Frederico Knabben. All rights reserved.
 * For licensing, see https://ckeditor.com/legal/ckeditor-oss-license
 */

CKEDITOR.editorConfig = function( config ) {
	// Define changes to default configuration here. For example:
	config.language = 'zh-cn';
    config.uiColor = '#d1d1d1';
    config.image_previewText = ' ';/*去掉图片预览框的文字*/
    /*开启工具栏“图像”中文件上传功能，后面的url为图片上传要指向的的action或servlet*/
    config.filebrowserImageUploadUrl= "/admin/sign/uploader_edit";
    config.fillEmptyBlocks  = false
};
