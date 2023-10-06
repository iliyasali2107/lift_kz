export default defineNuxtRouteMiddleware((to, from) => {
    // setTimeout(() => {
    //     if (process.server) return;
    //     console.log(localStorage.getItem('iin'));
    //     if (to.path !== '/auth/login' && !localStorage.getItem('iin')) {
    //         return navigateTo('/auth/login');
    //     }
    // });
});
