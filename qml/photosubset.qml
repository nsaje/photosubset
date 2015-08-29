import QtQuick 2.5
import QtQuick.Window 2.0

Rectangle {
    width: Screen.width
    height: Screen.height
    focus: true

		Keys.onLeftPressed: {
			controller.prev();
			event.accepted = true;
		}
		Keys.onRightPressed: {
			controller.next();
			event.accepted = true;
		}
		Keys.onPressed: {
			if (event.key >= 0x30 && event.key <= 0x39) {
        var num = event.key - 0x30;
				controller.tag(num);
        event.accepted = true;
			}
		}

		Image {
			id: imageContainer
			source: controller.currentPhotoPath
			fillMode: Image.PreserveAspectFit
      width: parent.width
      height: parent.height
      sourceSize.width: parent.width
      autoTransform: true
      anchors.horizontalCenter: parent.horizontalCenter
      anchors.verticalCenter: parent.verticalCenter
		}

    Text {
      id: tagsText
      text: controller.tagsText
      font.pointSize: 34
    }

    Binding {
      target: imageContainer
      property: "source"
      value: controller.currentPhotoPath
    }

    Binding {
      target: tagsText
      property: "text"
      value: controller.tagsText
    }
}
