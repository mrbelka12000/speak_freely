package linguo_sphere_backend

func GetFAQ(language string) string {
	translations := map[string]string{
		"en": en,
		"es": es,
		"fr": fr,
		"de": de,
		"it": it,
		"pt": pt,
		"ja": ja,
		"ko": ko,
		"ru": ru,
		"tr": tr,
	}

	val, ok := translations[language]
	if ok {
		return val
	}

	return translations[en]
}

var (
	en = `/start: Begins interaction with the service.
/language: Selects the language to practice.
/themes: Chooses a discussion theme.
/reset: Start a new general conversation. You can send an audio or video message and begin practicing.
/level: Sets the difficulty level of the topics.
/topic: Selects a specific topic to discuss.`

	es = `/start: Inicia la interacción con el servicio.
/language: Selecciona el idioma para practicar.
/themes: Elige un tema de conversación.
/reset: Inicia una nueva conversación general. Puedes enviar un mensaje de audio o video y comenzar a practicar.
/level: Establece el nivel de dificultad de los temas.
/topic: Selecciona un tema específico para discutir.`

	fr = `/start: Lance l'interaction avec le service.
/language: Sélectionne la langue à pratiquer.
/themes: Choisit un thème de discussion.
/reset: Inicia una nueva conversación general. Puedes enviar un mensaje de audio o video y comenzar a practicar.
/level: Définit le niveau de difficulté des sujets.
/topic: Sélectionne un sujet spécifique à discuter.`

	de = `/start: Startet die Interaktion mit dem Service.
/language: Wählt die Sprache zum Üben aus.
/themes: Wählt ein Diskussionsthema aus.
/reset: Starten Sie eine neue allgemeine Unterhaltung. Sie können eine Audio- oder Videonachricht senden und mit dem Üben beginnen.
/level: Legt das Schwierigkeitsniveau der Themen fest.
/topic: Wählt ein spezifisches Diskussionsthema aus.`

	it = `/start: Avvia l'interazione con il servizio.
/language: Seleziona la lingua da praticare.
/themes: Sceglie un tema di discussione.
/reset: Inizia una nuova conversazione generale. Puoi inviare un messaggio audio o video e iniziare a praticare.
/level: Imposta il livello di difficoltà degli argomenti.
/topic: Seleziona un argomento specifico da discutere.`

	pt = `/start: Inicia a interação com o serviço.
/language: Seleciona o idioma para praticar.
/themes: Escolhe um tema de discussão.
/reset: Inicie uma nova conversa geral. Você pode enviar uma mensagem de áudio ou vídeo e começar a praticar.
/level: Define o nível de dificuldade dos tópicos.
/topic: Seleciona um tópico específico para discutir.`

	ja = `/start: サービスとのやり取りを開始します。
/language: 練習する言語を選択します。
/themes: 議論するテーマを選択します。
/reset: 新しい一般的な会話を始めましょう。音声メッセージやビデオメッセージを送信して、練習を始めることができます。
/level: トピックの難易度レベルを設定します。
/topic: 議論する特定のトピックを選択します。`

	ko = `/start: 서비스와의 상호 작용을 시작합니다.
/language: 연습할 언어를 선택합니다.
/themes: 논의할 주제를 선택합니다.
/reset: 새로운 일반 대화를 시작하세요. 오디오 또는 비디오 메시지를 보내고 연습을 시작할 수 있습니다.
/level: 주제의 난이도 수준을 설정합니다.
/topic: 논의할 특정 주제를 선택합니다.`

	ru = `/start: Начинает взаимодействие с сервисом.
/language: Выбирает язык для практики.
/themes: Выбирает тему для обсуждения.
/reset: Начните новый общий разговор. Вы можете отправить аудио- или видеосообщение и начать практиковаться.
/level: Устанавливает уровень сложности тем.
/topic: Выбирает конкретную тему для обсуждения.`

	tr = `/start: Servis ile etkileşimi başlatır.
/language: Pratik yapmak için dili seçer.
/themes: Tartışılacak temayı seçer.
/reset: Yeni bir genel sohbet başlatın. Sesli veya görüntülü bir mesaj gönderebilir ve pratik yapmaya başlayabilirsiniz.
/level: Konuların zorluk seviyesini ayarlar.
/topic: Tartışılacak belirli bir konuyu seçer.`
)
