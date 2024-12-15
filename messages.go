package linguo_sphere_backend

func GetReachedLimitMessage(language string) string {
	translations := map[string]string{
		"en": "Sorry you have reached your limit",
		"es": "Lo siento, has alcanzado tu límite",
		"fr": "Désolé, vous avez atteint votre limite",
		"de": "Entschuldigung, Sie haben Ihr Limit erreicht",
		"it": "Spiacenti, hai raggiunto il tuo limite",
		"pt": "Desculpe, você atingiu seu limite",
		"ja": "申し訳ありませんが、上限に達しました",
		"ko": "죄송합니다. 한도에 도달했습니다",
		"ru": "Извините, вы достигли своего лимита",
		"tr": "Üzgünüz, sınırınıza ulaştınız",
	}

	val, ok := translations[language]
	if ok {
		return val
	}

	return translations["en"]
}

func GetTimeToPayMessage(language string) string {
	translations := map[string]string{
		"en": "We noticed your subscription has come to an end. Don't miss out — renew now to keep access to all the benefits you love!",
		"es": "Nos dimos cuenta de que tu suscripción ha llegado a su fin. ¡No te lo pierdas, renueva ahora para seguir disfrutando de todos los beneficios que te encantan!",
		"fr": "Nous avons remarqué que votre abonnement est arrivé à son terme. Ne manquez pas cette occasion — renouvelez maintenant pour continuer à profiter de tous les avantages que vous aimez !",
		"de": "Wir haben bemerkt, dass Ihr Abonnement abgelaufen ist. Verpassen Sie nichts — verlängern Sie jetzt, um weiterhin Zugriff auf alle Vorteile zu haben, die Sie lieben!",
		"it": "Abbiamo notato che il tuo abbonamento è terminato. Non perdere l'occasione — rinnova ora per continuare ad accedere a tutti i vantaggi che ami!",
		"pt": "Percebemos que sua assinatura chegou ao fim. Não perca — renove agora para continuar desfrutando de todos os benefícios que você ama!",
		"ja": "ご契約の期限が切れたことを確認しました。お見逃しなく — お気に入りのすべての特典に引き続きアクセスできるように、今すぐ更新してください！",
		"ko": "구독이 만료된 것을 확인했습니다. 놓치지 마세요 — 좋아하는 모든 혜택에 계속 액세스할 수 있도록 지금 갱신하세요!",
		"ru": "Мы заметили, что срок вашей подписки истёк. Не упустите возможность — продлите подписку сейчас, чтобы сохранить доступ ко всем любимым преимуществам!",
		"tr": "Aboneliğinizin sona erdiğini fark ettik. Kaçırmayın — sevdiğiniz tüm avantajlara erişimi sürdürmek için şimdi yenileyin!",
	}

	val, ok := translations[language]
	if ok {
		return val
	}

	return translations["en"]
}

func GetNothingFindMessage(language string) string {
	translations := map[string]string{
		"en": "Nothing was found for your criteria, try selecting a different topic.",
		"es": "No se encontró nada según tus criterios, intenta seleccionar un tema diferente.",
		"fr": "Rien n'a été trouvé selon vos critères, essayez de sélectionner un autre sujet.",
		"de": "Für Ihre Kriterien wurde nichts gefunden, versuchen Sie, ein anderes Thema auszuwählen.",
		"it": "Non è stato trovato nulla in base ai tuoi criteri, prova a selezionare un argomento diverso.",
		"pt": "Nada foi encontrado com seus critérios, tente selecionar um tema diferente.",
		"ja": "お客様の条件に一致するものは見つかりませんでした。別のトピックを選択してみてください。",
		"ko": "귀하의 기준에 맞는 항목을 찾을 수 없습니다. 다른 주제를 선택해 보세요.",
		"ru": "По вашим критериям ничего не найдено, попробуйте выбрать другую тему.",
		"tr": "Kriterlerinize uygun bir şey bulunamadı, farklı bir konu seçmeyi deneyin.",
	}

	val, ok := translations[language]
	if ok {
		return val
	}

	return translations["en"]
}

func GetAlreadyPaidMessage(language string) string {
	translations := map[string]string{
		"en": "It looks like you've already made a payment!",
		"es": "¡Parece que ya has realizado el pago!",
		"fr": "Il semble que vous ayez déjà effectué le paiement !",
		"de": "Es scheint, dass Sie bereits eine Zahlung vorgenommen haben!",
		"it": "Sembra che tu abbia già effettuato il pagamento!",
		"pt": "Parece que você já fez o pagamento!",
		"ja": "すでにお支払いが完了しているようです！",
		"ko": "이미 결제를 완료하신 것 같아요!",
		"ru": "Похоже, вы уже оплатили!",
		"tr": "Görünüşe göre zaten ödeme yapmışsınız!",
	}

	val, ok := translations[language]
	if ok {
		return val
	}

	return translations["en"]
}

func GetChooseLevelMessage(language string) string {
	translations := map[string]string{
		"en": "Which level are you?",
		"es": "¿En qué nivel estás?",
		"fr": "Quel est votre niveau ?",
		"de": "Auf welchem Niveau bist du?",
		"it": "A che livello sei?",
		"pt": "Em que nível você está?",
		"ja": "あなたのレベルはどのくらいですか？",
		"ko": "당신의 레벨은 무엇입니까?",
		"ru": "Какой у вас уровень?",
		"tr": "Hangi seviyedesiniz?",
	}

	val, ok := translations[language]
	if ok {
		return val
	}

	return translations["en"]
}

func GetChooseTopicMessage(language string) string {
	translations := map[string]string{
		"en": "Which topic do you want to discuss?",
		"es": "¿Qué tema quieres discutir?",
		"fr": "Quel sujet souhaitez-vous aborder ?",
		"de": "Welches Thema möchtest du besprechen?",
		"it": "Quale argomento vuoi discutere?",
		"pt": "Qual tópico você quer discutir?",
		"ja": "どのトピックについて話し合いたいですか？",
		"ko": "어떤 주제를 논의하고 싶으신가요?",
		"ru": "Какую тему вы хотите обсудить?",
		"tr": "Hangi konuyu tartışmak istersiniz?",
	}

	val, ok := translations[language]
	if ok {
		return val
	}

	return translations["en"]
}

func GetChooseThemeMessage(language string) string {
	translations := map[string]string{
		"en": "Which question do you want to discuss?",
		"es": "¿Qué pregunta quieres discutir?",
		"fr": "Quelle question souhaitez-vous aborder ?",
		"de": "Welche Frage möchtest du besprechen?",
		"it": "Quale domanda vuoi discutere?",
		"pt": "Qual pergunta você quer discutir?",
		"ja": "どの質問について話し合いたいですか？",
		"ko": "어떤 질문을 논의하고 싶으신가요?",
		"ru": "Какой вопрос вы хотите обсудить?",
		"tr": "Hangi soruyu tartışmak istersiniz?",
	}

	val, ok := translations[language]
	if ok {
		return val
	}

	return translations["en"]
}

func GetGreetingMessage(language string) string {
	translations := map[string]string{
		"en": "Hello %s, let's start to practice.",
		"es": "Hola %s, comencemos a practicar.",
		"fr": "Bonjour %s, commençons à pratiquer.",
		"de": "Hallo %s, lass uns mit dem Üben beginnen.",
		"it": "Ciao %s, iniziamo a praticare.",
		"pt": "Olá %s, vamos começar a praticar.",
		"ja": "こんにちは %s、練習を始めましょう。",
		"ko": "안녕하세요 %s, 연습을 시작해봅시다.",
		"ru": "Привет %s, давайте начнем практиковаться.",
		"tr": "Merhaba %s, haydi pratiğe başlayalım.",
	}

	val, ok := translations[language]
	if ok {
		return val
	}

	return translations["en"]
}

func GetYourChooseMessage(language string) string {
	translations := map[string]string{
		"en": "Your choice: ",
		"es": "Tu elección: ",
		"fr": "Votre choix: ",
		"de": "Deine Wahl: ",
		"it": "La tua scelta: ",
		"pt": "Sua escolha: ",
		"ja": "あなたの選択：",
		"ko": "당신의 선택: ",
		"ru": "Ваш выбор: ",
		"tr": "Senin seçimin: ",
	}

	val, ok := translations[language]
	if ok {
		return val
	}

	return translations["en"]
}
