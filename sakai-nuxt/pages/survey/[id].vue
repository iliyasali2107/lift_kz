<template>
    <Toast></Toast>
    <div>
        <Card style="margin-bottom: 0px">
            <template #title> Добро пожаловать</template>
            <template #content v-if="survey">
                <div style="display: flex; flex-direction: row">
                    <div style="display: flex; flex-direction: column; margin-right: auto">
                        <span>Название опроса: {{ survey.name }}</span>
                        <span style="margin-top: 10px">Адресс: {{ survey.adress }}</span>
                    </div>
                    <Button label="История" icon="pi pi-history" />
                    <Button label="Выйти" icon="pi pi-sign-out" severity="danger" style="margin-left: 10px" @click="exit" />
                </div>
                <!-- <hr /> -->
                <div class="flex justify-content-center" style="padding-left: 0px; padding-right: 0px">
                    <SelectButton v-model="value" :options="options" aria-labelledby="basic" />
                </div>
            </template>
        </Card>
        <Card v-if="survey" v-for="(question, key) in survey.QuestionDescriptions" class="question_class">
            <template #title> Вопрос {{ key + 1 }} </template>
            <template #content>
                <p>{{ question.description }} {{ question.description.includes('?') ? '' : '?' }}</p>
                <div class="flex align-items-center" style="margin-bottom: 2dvh">
                    <Checkbox v-model="results[key]" inputId="ingredient1" name="value" value="1" :disabled="results[key] && results[key] != `` && results[key] != `1`" />
                    <label for="ingredient1" class="ml-2"> Да </label>
                </div>
                <div class="flex align-items-center" style="margin-bottom: 2dvh">
                    <Checkbox v-model="results[key]" inputId="ingredient2" name="value" value="2" :disabled="results[key] && results[key] != `` && results[key] != `2`" />
                    <label for="ingredient2" class="ml-2"> Нет </label>
                </div>
                <div class="flex align-items-center" style="margin-bottom: 2dvh">
                    <Checkbox v-model="results[key]" inputId="ingredient3" name="value" value="3" :disabled="results[key] && results[key] != `` && results[key] != `3`" />
                    <label for="ingredient3" class="ml-2"> Воздержусь </label>
                </div>
                <InlineMessage v-if="errorRow == key">обязательно надо ответить</InlineMessage>
                <div class="p-flex-grow-1"></div>
                <!-- This div will push the button to the bottom -->
                <div class="text-right">
                    <Button label="Отменить голос" style="margin-bottom: 0; width: 15dvh; height: 5dvh; font-size: small" @click="results[key] = null" />
                </div>
            </template>
        </Card>
        <div v-else>
            <Skeleton class="question_class" height="20dvh" borderRadius="16px"></Skeleton>
            <Skeleton class="question_class" height="20dvh" borderRadius="16px"></Skeleton>
            <Skeleton class="question_class" height="20dvh" borderRadius="16px"></Skeleton>
        </div>
        <div class="center-button" style="margin-top: 5dvh"><Button label="Отправить" icon="pi pi-send" :loading="loading" @click="sendResults" /></div>
    </div>
</template>
<script>
import { useMainStore } from '../../service/mainstore';
export default {
    async setup() {
        definePageMeta({
            layout: false
        });
        const route = useRoute();
        const nuxtApp = useNuxtApp();
        const id = route.params.id;
        return { nuxtApp, id };
    },
    data() {
        return {
            loading: false,
            survey: null,
            errorRow: null,
            amount_of_question: 0,
            results: [],
            iin: useMainStore().get_iin,
            username: useMainStore().get_username,
            value: 'Жилая помещения',
            options: ['Жилая помещения', 'Нежилая помещения'],
            user_id: useMainStore().get_user_id
        };
    },

    async mounted() {
        await this.init();
    },
    methods: {
        async init() {
            if (!this.iin) this.iin = localStorage.getItem('iin');
            if (!this.user_id) this.user_id = localStorage.getItem('user_id');
            console.log('user_id:', this.user_id);
            this.survey = await this.nuxtApp.$liftservice().get_survey(this.id);
            console.log('response:', this.survey);
        },
        async sendResults() {
            this.loading = true;
            const questions = [];
            for (var i = 0; i < this.survey.QuestionDescriptions.length; i++) {
                try {
                    // console.log('QuestionDescriptions[i].description:', this.QuestionDescriptions[i].description);
                    questions[i] = {
                        id: this.survey.QuestionDescriptions[i].id,
                        answer_id: parseInt(this.results[i][0])
                    };
                } catch (e) {
                    this.errorRow = i;
                    this.loading = false;
                    return;
                }
            }
            const request = { id: parseInt(this.id), questions: questions, user_id: this.user_id };
            console.log('results:', request);
            try {
                await this.nuxtApp.$liftservice().post_answers(request);
                this.$toast.add({
                    severity: 'success',
                    summary: 'Успешно',
                    detail: 'Ваш голос отправлен.',
                    life: 3000
                });
            } catch (e) {
                this.$toast.add({
                    severity: 'error',
                    summary: 'Error',
                    detail: 'For unknown reason.',
                    life: 3000
                });
            }
            this.loading = false;
            this.clearResult();
        },
        clearResult() {
            this.results = [];
            this.errorRow = null;
        },
        exit() {
            useMainStore().clear_store();
            this.$router.push('/login/auth');
        }
    }
};
</script>
<style scoped>
/* Default styles for the Card element with margin-left and margin-right */
.question_class {
    margin-left: 15dvh;
    margin-right: 15dvh;
    margin-top: 5dvh;
}

/* Media query to set margin to 0 when the screen width matches the device screen */
@media (max-width: 600px) {
    .question_class {
        margin-left: 0;
        margin-right: 0;
    }
}
.center-button {
    display: flex;
    justify-content: center;
}
.p-button {
    margin-bottom: 50px;
}

@media (max-width: 600px) {
    .adaptive-select-button {
        font-size: medium; /* Change font size for smaller screens */
    }
}

@media (max-width: 400px) {
    .adaptive-select-button {
        font-size: small; /* Change font size for even smaller screens */
    }
}
</style>
