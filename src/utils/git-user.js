const exec = require( 'child_process' ).execSync;

module.exports = function() {
    return `${getName()} ${getEmail()}`.trim();
};

function getName() {
    try {
        return exec( 'git config --get user.name' ).toString().trim();
    } catch ( e ) {
        return '';
    }
}

function getEmail() {
    try {
        return `<${exec( 'git config --get user.email' ).toString().trim()}>`;
    } catch ( e ) {
        return '';
    }
}