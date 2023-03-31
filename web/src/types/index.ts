export interface Achievement {
	name: string;
	displayName: string;
	description: string;
	icon: string;

  achieved: boolean;
	unlockedAt: number;
}

export interface State {
  achievements: Achievement[];
}
