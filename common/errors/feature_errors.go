package errors

// AIROUTER: предупреждения о deprecated-фичах отключены. Ноды из подписок
// используют WS/VMess/Trojan/Shadowsocks/gRPC, переписать их нельзя, а спам
// в логе на каждый запуск не нужен. Сигнатуры оставлены, чтобы все места
// вызова компилировались без правок.
func PrintNonRemovalDeprecatedFeatureWarning(sourceFeature string, targetFeature string) {
}

// PrintDeprecatedFeatureWarning prints a warning for deprecated and going to be removed feature.
// Do not remove this function even there is no reference to it.
func PrintDeprecatedFeatureWarning(feature string, migrateFeature string) {
}

// PrintRemovedFeatureError prints an error message for removed feature then return an error. And after long enough time the message can also be removed, uses as an indicator.
// Do not remove this function even there is no reference to it.
func PrintRemovedFeatureError(feature string, migrateFeature string) error {
	if len(migrateFeature) > 0 {
		return New("The feature " + feature + " has been removed and migrated to " + migrateFeature + ". Please update your config(s) according to release note and documentation.")
	} else {
		return New("The feature " + feature + " has been removed. Please update your config(s) according to release note and documentation.")
	}
}
