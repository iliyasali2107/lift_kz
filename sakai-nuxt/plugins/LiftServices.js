import { API_URL, getHeader } from '@/config';

const LiftService = {
    async login() {
        const url = `${API_URL}/api/users/login`;
        const response = await useAsyncData('login', async () => await $fetch(url, { method: 'POST', headers: getHeader() }));

        return response;
    },
    async disableSurvey(survey_id) {
        const url = `${API_URL}/api/survey/close`;
        await useFetch(url, {
            method: 'POST',
            headers: getHeader(),
            body: { survey_id: survey_id }
        });
    },
    async getSurveyByID(survey_id) {
        const url = `${API_URL}/api/survey/summary/${survey_id}`;
        const response = await useFetch(url, { headers: getHeader() });
        return response;
    },
    async confirm(requirements) {
        const url = `${API_URL}/api/users/confirm`;
        const response = await useAsyncData('confirm', async () => await $fetch(url, { method: 'POST', headers: getHeader(), body: requirements }));
        return response;
    },
    async get_surveys() {
        const url = `${API_URL}/api/survey/get/surveys/1`;
        console.log('URL:', url);
        // const data = await useAsyncData('survey', async () => await $fetch(() => url, { method: 'GET', headers: getHeader() }));
        const data = await $fetch(url);
        console.log('GXDZSFGHJ');
        // const response = await useFetch(`${API_URL}/api/mams`);
        return data;
    },
    async post_survey(requirements) {
        const url = `${API_URL}/api/survey/create`;
        const response = await useAsyncData('createproduct', async () => await $fetch(url, { method: 'POST', headers: getHeader(), body: requirements }));
        return response;
    },
    async get_survey(survey_id) {
        const url = `${API_URL}/api/survey/get/survey/${survey_id}`;
        const response = await $fetch(url, { headers: getHeader() });
        return response;
    }
};
export default defineNuxtPlugin((nuxtApp) => {
    // Doing something with nuxtApp
    nuxtApp.provide('liftservice', () => LiftService);
});
