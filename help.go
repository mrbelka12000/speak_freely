package linguo_sphere_backend

func GetFAQ(language string) string {
	mp := map[string]string{
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

	val, ok := mp[language]
	if ok {
		return val
	}

	return mp[en]
}

var (
	en = `/start: Begins interaction with the service.
/language: Selects the language to practice.
/themes: Chooses a discussion theme.
/conversation: Starts a new general conversation and deletes the previous conversation from memory.
/level: Sets the difficulty level of the topics.
/topic: Selects a specific topic to discuss.`

	es = `/start: Inicia la interacción con el servicio.
/language: Selecciona el idioma para practicar.
/themes: Elige un tema de conversación.
/conversation: Inicia una nueva conversación general y elimina la conversación anterior de la memoria.
/level: Establece el nivel de dificultad de los temas.
/topic: Selecciona un tema específico para discutir.`

	fr = `/start: Lance l'interaction avec le service.
/language: Sélectionne la langue à pratiquer.
/themes: Choisit un thème de discussion.
/conversation: Démarre une nouvelle conversation générale et supprime la conversation précédente de la mémoire.
/level: Définit le niveau de difficulté des sujets.
/topic: Sélectionne un sujet spécifique à discuter.`

	de = `/start: Startet die Interaktion mit dem Service.
/language: Wählt die Sprache zum Üben aus.
/themes: Wählt ein Diskussionsthema aus.
/conversation: Startet ein neues allgemeines Gespräch und löscht das vorherige Gespräch aus dem Speicher.
/level: Legt das Schwierigkeitsniveau der Themen fest.
/topic: Wählt ein spezifisches Diskussionsthema aus.`

	it = `/start: Avvia l'interazione con il servizio.
/language: Seleziona la lingua da praticare.
/themes: Sceglie un tema di discussione.
/conversation: Avvia una nuova conversazione generale e elimina la conversazione precedente dalla memoria.
/level: Imposta il livello di difficoltà degli argomenti.
/topic: Seleziona un argomento specifico da discutere.`

	pt = `/start: Inicia a interação com o serviço.
/language: Seleciona o idioma para praticar.
/themes: Escolhe um tema de discussão.
/conversation: Inicia uma nova conversa geral e apaga a conversa anterior da memória.
/level: Define o nível de dificuldade dos tópicos.
/topic: Seleciona um tópico específico para discutir.`

	ja = `/start: サービスとのやり取りを開始します。
/language: 練習する言語を選択します。
/themes: 議論するテーマを選択します。
/conversation: 新しい一般的な会話を開始し、前の会話をメモリから削除します。
/level: トピックの難易度レベルを設定します。
/topic: 議論する特定のトピックを選択します。`

	ko = `/start: 서비스와의 상호 작용을 시작합니다.
/language: 연습할 언어를 선택합니다.
/themes: 논의할 주제를 선택합니다.
/conversation: 새 일반 대화를 시작하고 이전 대화를 메모리에서 삭제합니다.
/level: 주제의 난이도 수준을 설정합니다.
/topic: 논의할 특정 주제를 선택합니다.`

	ru = `/start: Начинает взаимодействие с сервисом.
/language: Выбирает язык для практики.
/themes: Выбирает тему для обсуждения.
/conversation: Начинает новый общий разговор и удаляет предыдущий разговор из памяти.
/level: Устанавливает уровень сложности тем.
/topic: Выбирает конкретную тему для обсуждения.`

	tr = `/start: Servis ile etkileşimi başlatır.
/language: Pratik yapmak için dili seçer.
/themes: Tartışılacak temayı seçer.
/conversation: Yeni bir genel konuşma başlatır ve önceki konuşmayı hafızadan siler.
/level: Konuların zorluk seviyesini ayarlar.
/topic: Tartışılacak belirli bir konuyu seçer.`
)
