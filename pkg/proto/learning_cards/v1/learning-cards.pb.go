// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.27.1
// source: learning_cards/v1/learning-cards.proto

package learning_cards

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_learning_cards_v1_learning_cards_proto protoreflect.FileDescriptor

const file_learning_cards_v1_learning_cards_proto_rawDesc = "" +
	"\n" +
	"&learning_cards/v1/learning-cards.proto\x12\x11learning_cards.v1\x1a\x1cgoogle/api/annotations.proto\x1a\x1bgoogle/protobuf/empty.proto\x1a\x1clearning_cards/v1/card.proto\x1a\x1dlearning_cards/v1/group.proto2\xe3\b\n" +
	"\rLearningCards\x12r\n" +
	"\vCreateGroup\x12%.learning_cards.v1.CreateGroupRequest\x1a&.learning_cards.v1.CreateGroupResponse\"\x14\x82\xd3\xe4\x93\x02\x0e:\x01*\"\t/v1/group\x12l\n" +
	"\n" +
	"ListGroups\x12$.learning_cards.v1.ListGroupsRequest\x1a%.learning_cards.v1.ListGroupsResponse\"\x11\x82\xd3\xe4\x93\x02\v\x12\t/v1/group\x12q\n" +
	"\bGetGroup\x12\".learning_cards.v1.GetGroupRequest\x1a#.learning_cards.v1.GetGroupResponse\"\x1c\x82\xd3\xe4\x93\x02\x16\x12\x14/v1/group/{group_id}\x12m\n" +
	"\vUpdateGroup\x12%.learning_cards.v1.UpdateGroupRequest\x1a\x16.google.protobuf.Empty\"\x1f\x82\xd3\xe4\x93\x02\x19:\x01*\x1a\x14/v1/group/{group_id}\x12j\n" +
	"\vDeleteGroup\x12%.learning_cards.v1.DeleteGroupRequest\x1a\x16.google.protobuf.Empty\"\x1c\x82\xd3\xe4\x93\x02\x16*\x14/v1/group/{group_id}\x12e\n" +
	"\aAddCard\x12!.learning_cards.v1.AddCardRequest\x1a\".learning_cards.v1.AddCardResponse\"\x13\x82\xd3\xe4\x93\x02\r:\x01*\"\b/v1/card\x12z\n" +
	"\tListCards\x12#.learning_cards.v1.ListCardsRequest\x1a$.learning_cards.v1.ListCardsResponse\"\"\x82\xd3\xe4\x93\x02\x1c\x12\x1a/v1/group/{group_id}/cards\x12l\n" +
	"\aGetCard\x12!.learning_cards.v1.GetCardRequest\x1a\".learning_cards.v1.GetCardResponse\"\x1a\x82\xd3\xe4\x93\x02\x14\x12\x12/v1/card/{card_id}\x12i\n" +
	"\n" +
	"UpdateCard\x12$.learning_cards.v1.UpdateCardRequest\x1a\x16.google.protobuf.Empty\"\x1d\x82\xd3\xe4\x93\x02\x17:\x01*\x1a\x12/v1/card/{card_id}\x12f\n" +
	"\n" +
	"DeleteCard\x12$.learning_cards.v1.DeleteCardRequest\x1a\x16.google.protobuf.Empty\"\x1a\x82\xd3\xe4\x93\x02\x14*\x12/v1/card/{card_id}BRZPgithub.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1;learning_cardsb\x06proto3"

var file_learning_cards_v1_learning_cards_proto_goTypes = []any{
	(*CreateGroupRequest)(nil),  // 0: learning_cards.v1.CreateGroupRequest
	(*ListGroupsRequest)(nil),   // 1: learning_cards.v1.ListGroupsRequest
	(*GetGroupRequest)(nil),     // 2: learning_cards.v1.GetGroupRequest
	(*UpdateGroupRequest)(nil),  // 3: learning_cards.v1.UpdateGroupRequest
	(*DeleteGroupRequest)(nil),  // 4: learning_cards.v1.DeleteGroupRequest
	(*AddCardRequest)(nil),      // 5: learning_cards.v1.AddCardRequest
	(*ListCardsRequest)(nil),    // 6: learning_cards.v1.ListCardsRequest
	(*GetCardRequest)(nil),      // 7: learning_cards.v1.GetCardRequest
	(*UpdateCardRequest)(nil),   // 8: learning_cards.v1.UpdateCardRequest
	(*DeleteCardRequest)(nil),   // 9: learning_cards.v1.DeleteCardRequest
	(*CreateGroupResponse)(nil), // 10: learning_cards.v1.CreateGroupResponse
	(*ListGroupsResponse)(nil),  // 11: learning_cards.v1.ListGroupsResponse
	(*GetGroupResponse)(nil),    // 12: learning_cards.v1.GetGroupResponse
	(*emptypb.Empty)(nil),       // 13: google.protobuf.Empty
	(*AddCardResponse)(nil),     // 14: learning_cards.v1.AddCardResponse
	(*ListCardsResponse)(nil),   // 15: learning_cards.v1.ListCardsResponse
	(*GetCardResponse)(nil),     // 16: learning_cards.v1.GetCardResponse
}
var file_learning_cards_v1_learning_cards_proto_depIdxs = []int32{
	0,  // 0: learning_cards.v1.LearningCards.CreateGroup:input_type -> learning_cards.v1.CreateGroupRequest
	1,  // 1: learning_cards.v1.LearningCards.ListGroups:input_type -> learning_cards.v1.ListGroupsRequest
	2,  // 2: learning_cards.v1.LearningCards.GetGroup:input_type -> learning_cards.v1.GetGroupRequest
	3,  // 3: learning_cards.v1.LearningCards.UpdateGroup:input_type -> learning_cards.v1.UpdateGroupRequest
	4,  // 4: learning_cards.v1.LearningCards.DeleteGroup:input_type -> learning_cards.v1.DeleteGroupRequest
	5,  // 5: learning_cards.v1.LearningCards.AddCard:input_type -> learning_cards.v1.AddCardRequest
	6,  // 6: learning_cards.v1.LearningCards.ListCards:input_type -> learning_cards.v1.ListCardsRequest
	7,  // 7: learning_cards.v1.LearningCards.GetCard:input_type -> learning_cards.v1.GetCardRequest
	8,  // 8: learning_cards.v1.LearningCards.UpdateCard:input_type -> learning_cards.v1.UpdateCardRequest
	9,  // 9: learning_cards.v1.LearningCards.DeleteCard:input_type -> learning_cards.v1.DeleteCardRequest
	10, // 10: learning_cards.v1.LearningCards.CreateGroup:output_type -> learning_cards.v1.CreateGroupResponse
	11, // 11: learning_cards.v1.LearningCards.ListGroups:output_type -> learning_cards.v1.ListGroupsResponse
	12, // 12: learning_cards.v1.LearningCards.GetGroup:output_type -> learning_cards.v1.GetGroupResponse
	13, // 13: learning_cards.v1.LearningCards.UpdateGroup:output_type -> google.protobuf.Empty
	13, // 14: learning_cards.v1.LearningCards.DeleteGroup:output_type -> google.protobuf.Empty
	14, // 15: learning_cards.v1.LearningCards.AddCard:output_type -> learning_cards.v1.AddCardResponse
	15, // 16: learning_cards.v1.LearningCards.ListCards:output_type -> learning_cards.v1.ListCardsResponse
	16, // 17: learning_cards.v1.LearningCards.GetCard:output_type -> learning_cards.v1.GetCardResponse
	13, // 18: learning_cards.v1.LearningCards.UpdateCard:output_type -> google.protobuf.Empty
	13, // 19: learning_cards.v1.LearningCards.DeleteCard:output_type -> google.protobuf.Empty
	10, // [10:20] is the sub-list for method output_type
	0,  // [0:10] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_learning_cards_v1_learning_cards_proto_init() }
func file_learning_cards_v1_learning_cards_proto_init() {
	if File_learning_cards_v1_learning_cards_proto != nil {
		return
	}
	file_learning_cards_v1_card_proto_init()
	file_learning_cards_v1_group_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_learning_cards_v1_learning_cards_proto_rawDesc), len(file_learning_cards_v1_learning_cards_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_learning_cards_v1_learning_cards_proto_goTypes,
		DependencyIndexes: file_learning_cards_v1_learning_cards_proto_depIdxs,
	}.Build()
	File_learning_cards_v1_learning_cards_proto = out.File
	file_learning_cards_v1_learning_cards_proto_goTypes = nil
	file_learning_cards_v1_learning_cards_proto_depIdxs = nil
}
