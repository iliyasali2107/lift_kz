    CREATE TABLE public.answer (
        id integer NOT NULL,
        name character varying(100)
    );


   


  


    CREATE TABLE public.user_question (
        id integer NOT NULL,
        user_id integer NOT NULL,
        question_id integer NOT NULL,
        answer_id integer
    );