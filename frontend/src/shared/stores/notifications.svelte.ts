class NotificationsStore {
	hydrated = $state(false);

	unreadCount = $state(0);

	seed(count: number) {
		this.unreadCount = count;
    
		this.hydrated = true;
	}

	markOneSeen() {
		if (this.unreadCount > 0) this.unreadCount -= 1;
	}
}

export const notificationsStore = new NotificationsStore();
