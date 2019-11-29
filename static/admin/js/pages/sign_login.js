/*
 *  Document   : op_auth_signin.js
 *  Author     : pixelcave
 *  Description: Custom JS code used in Sign In Page
 */
class pageAuthSignIn {
    static initValidation() {
        // Load default options for jQuery Validation plugin
        One.helpers('validation');
        // Init Form Validation
        jQuery('.js-validation-signin').validate({
            rules: {
                'login-username': {
                    required: true,
                    minlength: 3
                },
                'login-password': {
                    required: true,
                    minlength: 5
                }
            },
            messages: {
                'login-username': {
                    required: '请输入用户名',
                    minlength: '您的用户名必须包含至少3个字符'
                },
                'login-password': {
                    required: '请输入密码',
                    minlength: '您的密码必须至少有5个字符长'
                }
            },
            submitHandler: function(form){
               app.post($(form).attr('action'),$(form).serialize());
            }
        });
    }
    static init() {
        this.initValidation();
    }
}
// Initialize when page loads
jQuery(() => { pageAuthSignIn.init(); });
