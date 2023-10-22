<template>
    <div>
        <Card style="margin-bottom: 0px" v-if="iin">
            <template #title> Добро пожаловать, {{ username }} </template>
            <template #content v-if="survey">
                <div style="display: flex; flex-direction: row">
                    <div style="display: flex; flex-direction: column; margin-right: auto">
                        <span>Название опроса: {{ survey.name }}</span>
                        <span style="margin-top: 10px">Адресс: {{ survey.adress }}</span>
                    </div>
                    <Button label="История" icon="pi pi-history" />
                </div>
            </template>
        </Card>
        <div class="card flex justify-content-center" style="padding-left: 0px; padding-right: 0px">
            <SelectButton v-model="value" :options="options" aria-labelledby="basic" />
        </div>
        <Card v-if="survey" v-for="(question, key) in survey.questions_str" class="question_class">
            <template #title> Вопрос {{ key + 1 }} </template>
            <template #content>
                <p>{{ question }} {{ question.includes('?') ? '' : '?' }}</p>
                <div class="flex align-items-center" style="margin-bottom: 2dvh">
                    <Checkbox v-model="results[key]" inputId="ingredient1" name="value" value="Да" />
                    <label for="ingredient1" class="ml-2"> Да </label>
                </div>
                <div class="flex align-items-center" style="margin-bottom: 2dvh">
                    <Checkbox v-model="results[key]" inputId="ingredient2" name="value" value="Нет" />
                    <label for="ingredient2" class="ml-2"> Нет </label>
                </div>
                <div class="flex align-items-center">
                    <Checkbox v-model="results[key]" inputId="ingredient3" name="value" value="Воздержусь" />
                    <label for="ingredient3" class="ml-2"> Воздержусь </label>
                </div>
            </template>
        </Card>
        <div v-else>
            <Skeleton class="question_class" height="20dvh" borderRadius="16px"></Skeleton>
            <Skeleton class="question_class" height="20dvh" borderRadius="16px"></Skeleton>
            <Skeleton class="question_class" height="20dvh" borderRadius="16px"></Skeleton>
        </div>
        <div class="center-button" style="margin-top: 5dvh"><Button label="Отправить" icon="pi pi-send" @click="sendResults" /></div>
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
        return { survey: null, amount_of_question: 0, results: [], iin: useMainStore().get_iin, username: useMainStore().get_username, value: 'Жилая помещения', options: ['Жилая помещения', 'Нежилая помещения'] };
    },

    async mounted() {
        await this.init();
    },
    methods: {
        async init() {
            this.survey = await this.nuxtApp.$liftservice().get_survey(this.id);
            console.log('response:', this.survey);
        },
        sendResults() {
            console.log('results:', this.results);
        }
        // optionSelected(key, value) {
        //     this.survey['answer'][key] = value;
        //     console.log(this.survey);
        // }
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
