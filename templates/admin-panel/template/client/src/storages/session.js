import createStorage from 'sham-ui-data-storage';

export const { useStorage } = createStorage( {
    name: '',
    email: '',
    sessionValidated: false,
    isAuthenticated: false
}, {
    DI: 'session:storage'
} );
